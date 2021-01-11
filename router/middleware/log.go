package middleware

import (
	"avito_mx/config"
	"net/http"
)

func Log(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		config.Logger.Println(r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}
