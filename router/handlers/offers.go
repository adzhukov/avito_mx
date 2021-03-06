package handlers

import (
	"avito_mx/config"
	"avito_mx/models"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func OffersHandler(w http.ResponseWriter, r *http.Request) {
	urlQuery := r.URL.Query()
	query := "SELECT offer_id, seller_id, name, price, quantity FROM offers"

	filtersQuery := make([]string, 0, 3)
	filterValues := make([]interface{}, 0, 3)

	seller := urlQuery.Get("seller_id")
	if seller != "" {
		sellerID, err := strconv.ParseInt(seller, 10, 64)
		if err != nil {
			responseJSON(w, respError{"parameter seller_id is not valid integer"}, http.StatusBadRequest)
			return
		}
		filterValues = append(filterValues, sellerID)
		filtersQuery = append(filtersQuery, fmt.Sprintf("seller_id = $%d", len(filterValues)))
	}

	offer := urlQuery.Get("offer_id")
	if offer != "" {
		offerID, err := strconv.ParseInt(offer, 10, 64)
		if err != nil {
			responseJSON(w, respError{"parameter offer_id is not valid integer"}, http.StatusBadRequest)
			return
		}
		filterValues = append(filterValues, offerID)
		filtersQuery = append(filtersQuery, fmt.Sprintf("offer_id = $%d", len(filterValues)))
	}

	q := urlQuery.Get("q")
	if q != "" {
		filterValues = append(filterValues, q)
		filtersQuery = append(filtersQuery, fmt.Sprintf("name ILIKE '%%' || $%d || '%%'", len(filterValues)))
	}

	if len(filtersQuery) > 0 {
		query += " WHERE " + strings.Join(filtersQuery, " AND ")
	}

	rows, err := config.DB.Query(r.Context(), query, filterValues...)
	if err != nil {
		config.Logger.Println("Unable conn.Query", err)
		responseJSON(w, respError{"A database error"}, http.StatusInternalServerError)
		return
	}

	offers := make([]models.Offer, 0)
	var row models.Offer

	for rows.Next() {
		err := rows.Scan(&row.OfferID, &row.SellerID, &row.Name, &row.Price, &row.Quantity)
		if err != nil {
			config.Logger.Println(err)
			break
		}
		offers = append(offers, row)
	}

	responseJSON(w, offers, http.StatusOK)
}
