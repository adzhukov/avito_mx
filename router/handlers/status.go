package handlers

import (
	"avito_mx/controllers"
	"fmt"
	"net/http"
	"strconv"
)

func StatusHandler(w http.ResponseWriter, r *http.Request) {
	taskStr := r.URL.Query().Get("task_id")
	if taskStr == "" {
		responseJSON(w, respError{"task_id parameter required"}, http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(taskStr, 10, 64)
	if err != nil {
		responseJSON(w, respError{"value of task_id parameter is invalid"}, http.StatusBadRequest)
		return
	}

	task, err := controllers.GetTaskByID(r.Context(), id)
	if err != nil {
		responseJSON(w, respError{fmt.Sprint(err)}, http.StatusBadRequest)
		return
	}

	responseJSON(w, task, http.StatusOK)
}
