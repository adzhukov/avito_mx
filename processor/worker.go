package processor

import (
	"avito_mx/config"
	"avito_mx/controllers"
	"avito_mx/models"
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/tealeg/xlsx/v3"
)

type parsedXlsx struct {
	Invalid   int
	Available []models.Offer
	Deleted   []int
}

const taskTimeout = 30

const (
	offerIDColumn   = iota
	nameColumn      = iota
	priceColumn     = iota
	quantityColumn  = iota
	availableColumn = iota
)

var client http.Client

func init() {
	t := &http.Transport{
		Dial: (&net.Dialer{
			Timeout:   15 * time.Second,
			KeepAlive: 15 * time.Second,
		}).Dial,
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: 5 * time.Second,
	}

	if os.Getenv("DEBUG") != "" {
		t.RegisterProtocol("file", http.NewFileTransport(http.Dir("/")))
	}

	client = http.Client{
		Transport: t,
	}
}

func (p *Processor) Worker() {
	for task := range p.queue {
		ctx, cancel := context.WithTimeout(context.Background(), taskTimeout*time.Second)
		ProcessTask(ctx, &task)
		cancel()
	}

	p.wg.Done()
}

func ProcessTask(ctx context.Context, task *models.Task) {
	taskReceived := time.Now()
	err := processTask(ctx, task)

	elapsed := time.Since(taskReceived)

	if err != nil {
		task.Status = models.TaskFailed
		task.Error = fmt.Sprint(err)
	} else {
		task.Status = models.TaskSuccess
	}

	config.Logger.Printf("Task %d processed in %s with result %s", task.TaskID, elapsed, task.Status)
	err = controllers.UpdateTask(task)
	if err != nil {
		config.Logger.Println(err)
	}
}

func processTask(ctx context.Context, task *models.Task) error {
	parsed, err := parseTask(task)
	if err != nil {
		return err
	}

	stats, err := send(ctx, task.SellerID, parsed)
	if err != nil {
		return err
	}

	task.Stats = stats
	return nil
}

func parseTask(task *models.Task) (*parsedXlsx, error) {
	response, err := client.Get(task.FileURL)
	if err != nil {
		config.Logger.Println("Unable to get file", err)
		return nil, err
	}

	file, err := ioutil.ReadAll(response.Body)
	if err != nil {
		config.Logger.Println("Unable to read response", err)
		return nil, err
	}

	wb, err := xlsx.OpenBinary(file)
	if err != nil {
		config.Logger.Println("Unable to parse file", err)
		return nil, err
	}

	var parsed parsedXlsx
	for _, sheet := range wb.Sheets {
		parseSheet(sheet, &parsed)
	}

	return &parsed, nil
}

func parseSheet(sheet *xlsx.Sheet, parsed *parsedXlsx) {
	sheet.ForEachRow(func(row *xlsx.Row) error {
		offerID, err := row.GetCell(offerIDColumn).Int()
		if err != nil {
			parsed.Invalid++
			return nil
		}

		available, err := strconv.ParseBool(row.GetCell(availableColumn).String())
		if err != nil {
			parsed.Invalid++
			return nil
		}

		quantity, err := row.GetCell(quantityColumn).Int()
		if err != nil || quantity < 0 {
			parsed.Invalid++
			return nil
		}

		if !available || quantity == 0 {
			parsed.Deleted = append(parsed.Deleted, offerID)
			return nil
		}

		price, err := row.GetCell(priceColumn).Int()
		if err != nil || price < 0 {
			parsed.Invalid++
			return nil
		}

		name := row.GetCell(nameColumn).Value
		if name == "" {
			parsed.Invalid++
			return nil
		}

		parsed.Available = append(parsed.Available,
			models.Offer{
				OfferID:  offerID,
				Name:     name,
				Price:    price,
				Quantity: quantity,
			})

		return nil
	})
}

func send(ctx context.Context, seller int64, offers *parsedXlsx) (*models.TaskStats, error) {
	tx, err := config.DB.Begin(ctx)
	if err != nil {
		return nil, err
	}

	batch := &pgx.Batch{}
	query := `INSERT INTO offers(seller_id, offer_id, name, price, quantity)
	VALUES ($1, $2, $3, $4, $5)
	ON CONFLICT (seller_id, offer_id)
	DO UPDATE SET name = EXCLUDED.name, price = EXCLUDED.price, quantity = EXCLUDED.quantity
	RETURNING (NOT xmax = 0)`

	for _, offer := range offers.Available {
		batch.Queue(query, seller, offer.OfferID, offer.Name, offer.Price, offer.Quantity)
	}

	query = "DELETE FROM offers WHERE seller_id = $1 AND offer_id = $2"

	for _, offer := range offers.Deleted {
		batch.Queue(query, seller, offer)
	}

	batchResults := tx.SendBatch(ctx, batch)

	stats := models.TaskStats{
		Invalid: offers.Invalid,
	}

	for range offers.Available {
		rows, err := batchResults.Query()
		if err != nil {
			return nil, err
		}

		for rows.Next() {
			var updated bool
			err := rows.Scan(&updated)
			if err != nil {
				return nil, err
			}
			if updated {
				stats.Updated++
			} else {
				stats.Created++
			}
		}
	}

	for range offers.Deleted {
		ct, err := batchResults.Exec()
		if err != nil {
			return nil, err
		}
		stats.Deleted += int(ct.RowsAffected())
	}

	batchResults.Close()
	return &stats, tx.Commit(ctx)
}
