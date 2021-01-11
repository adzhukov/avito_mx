package router

import (
	"avito_mx/router/handlers"
	"avito_mx/router/middleware"
	"net/http"
)

func New() http.Handler {
	router := http.NewServeMux()

	status := http.HandlerFunc(handlers.StatusHandler)
	offers := http.HandlerFunc(handlers.OffersHandler)
	importAsync := http.HandlerFunc(handlers.ImportHandler)
	importSync := http.HandlerFunc(handlers.SyncImportHandler)

	router.Handle("/status", middleware.Log(status))
	router.Handle("/offers", middleware.Log(offers))
	router.Handle("/import", middleware.Log(importAsync))
	router.Handle("/import/sync", middleware.Log(importSync))

	return router
}
