package handlers

import (
	"avito_mx/controllers"
	"avito_mx/models"
	"avito_mx/processor"
	"net/http"
	"net/url"
	"strconv"
)

func SyncImportHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	urlString := query.Get("url")
	if urlString == "" {
		responseJSON(w, respError{"parameter url is required"}, http.StatusBadRequest)
		return
	}
	_, err := url.ParseRequestURI(urlString)
	if err != nil {
		responseJSON(w, respError{"parameter url is not valid URI"}, http.StatusBadRequest)
		return
	}

	seller := query.Get("seller_id")
	if seller == "" {
		responseJSON(w, respError{"parameter seller_id is required"}, http.StatusBadRequest)
		return
	}
	sellerID, err := strconv.ParseInt(seller, 10, 64)
	if err != nil {
		responseJSON(w, respError{"parameter seller_id is not valid integer"}, http.StatusBadRequest)
		return
	}

	taskID, err := controllers.NewTask(r.Context(), sellerID, urlString)
	if err != nil {
		responseJSON(w, respError{"Unable to create task"}, http.StatusInternalServerError)
		return
	}

	task := models.Task{
		Status: models.TaskQueued,
		TaskID: taskID,
		TaskInfo: models.TaskInfo{
			SellerID: sellerID,
			FileURL:  urlString,
		},
	}

	processor.ProcessTask(r.Context(), &task)

	responseJSON(w, task, http.StatusOK)
}
