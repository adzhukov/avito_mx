package config

import (
	"avito_mx/models"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
)

var DB *pgxpool.Pool
var Logger *log.Logger
var Queue chan models.Task
