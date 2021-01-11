package handlers

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type respError struct {
	ErrorMessage string `json:"error"`
}

func responseJSON(w http.ResponseWriter, response interface{}, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(response)
}
