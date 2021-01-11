package main

import (
	"avito_mx/config"
	"avito_mx/processor"
	"avito_mx/router"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

const addr = ":8080"

func newDB() *pgxpool.Pool {
	pgConfig := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_DB"),
	)

	poolConfig, err := pgxpool.ParseConfig(pgConfig)
	if err != nil {
		log.Panic("Unable to parse postgres config", err)
	}

	poolConfig.MaxConns = int32(runtime.NumCPU() * 2)
	poolConfig.MinConns = int32(runtime.NumCPU())

	db, err := pgxpool.ConnectConfig(context.Background(), poolConfig)
	if err != nil {
		log.Panic("Unable to connect to postgres", err)
	}

	return db
}

func main() {
	config.Logger = log.New(os.Stdout, "", log.LstdFlags)
	config.DB = newDB()

	taskProcessor := processor.New()
	taskProcessor.Start()

	server := &http.Server{
		Addr:         addr,
		Handler:      router.New(),
		ErrorLog:     config.Logger,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	done := make(chan struct{}, 1)

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		config.Logger.Println("Got signal, shutdown")

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		server.SetKeepAlivesEnabled(false)
		if err := server.Shutdown(ctx); err != nil {
			config.Logger.Fatalf("Could not gracefully shutdown: %v\n", err)
		}

		close(config.Queue)
		taskProcessor.Wait()

		close(done)
	}()

	config.Logger.Println("Starting server at", addr)

	go func() {
		if err := server.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				config.Logger.Fatalf("Unexpected server closing: %v", err)
			}
		}
	}()

	<-done
}
