package main

import (
	"log/slog"
	"os"

	"github.com/alcb1310/final-bca-go/internal/database"
	"github.com/alcb1310/final-bca-go/internal/router"
	_ "github.com/joho/godotenv/autoload"
)

var port = os.Getenv("PORT")

func main() {
	connStr := os.Getenv("DATABASE_URL")
	db, data := database.New(connStr)
	if db == nil {
		slog.Error("New Database: Unable to connect to database")
		os.Exit(1)
	}

	if err := database.CreateTables(data); err != nil {
		slog.Error("New Database: Unable to create tables", "error", err)
		os.Exit(1)
	}

	r := router.NewRouter(db, port)
	if r == nil {
		os.Exit(1)
	}

	slog.Info("Server starting", "port", port)
	if err := r.ListenAndServe(); err != nil {
		slog.Error("Error starting server", "error", err)
		os.Exit(1)
	}
}
