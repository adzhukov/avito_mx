package controllers

import (
	"avito_mx/config"
	"avito_mx/models"
	"context"
	"errors"
)

func NewTask(ctx context.Context, seller int64, file string) (int64, error) {
	query := `INSERT INTO tasks (seller_id, file_url, status)
	VALUES ($1, $2, 'Queued')
	RETURNING task_id;`

	row := config.DB.QueryRow(ctx, query, seller, file)

	var id int64
	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func GetTaskByID(ctx context.Context, id int64) (*models.Task, error) {
	row := config.DB.QueryRow(ctx, "SELECT * from tasks where task_id=$1", id)
	var task models.Task
	task.Stats = new(models.TaskStats)

	err := row.Scan(&task.TaskID, &task.SellerID, &task.FileURL, &task.Status,
		&task.Stats.Created, &task.Stats.Updated, &task.Stats.Deleted, &task.Stats.Invalid)
	if err != nil {
		return nil, err
	}

	return &task, nil
}

func UpdateTask(task *models.Task) error {
	query := `UPDATE tasks SET status=$1`

	args := make([]interface{}, 2, 6)
	args[0] = task.Status
	args[1] = task.TaskID

	if task.Status == models.TaskSuccess {
		query += ", created=$3, updated=$4, deleted=$5, invalid=$6"
		args = append(args,
			task.Stats.Created,
			task.Stats.Updated,
			task.Stats.Deleted,
			task.Stats.Invalid)
	} else {
		query += ", error=$3"
		args = append(args, task.Error)
	}

	query += "WHERE task_id=$2"

	ct, err := config.DB.Exec(context.Background(), query, args...)
	if err != nil {
		return err
	}

	if ct.RowsAffected() == 0 {
		return errors.New("Task not found")
	}

	return nil
}
