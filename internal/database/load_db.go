package database

import (
	"database/sql"
	"log/slog"
	"os"
	"strings"
)

func CreateTables(d *sql.DB) error {
	file, err := os.ReadFile("./schema/tables.sql")
	if err != nil {
		slog.Error("Error reading file", "error", err)
		return err
	}

	tx, err := d.Begin()
	if err != nil {
		slog.Error("Error creating transaction", "error", err)
		return err
	}
	defer func() {
		if err := tx.Commit(); err != nil {
			slog.Error("Error committing transaction", "error", err)
		}
	}()

	queries := strings.SplitSeq(string(file), ";")
	for query := range queries {
		if _, err := tx.Exec(query); err != nil {
			slog.Error("Error executing query", "error", err, "query", query)
			tx.Rollback()
			return err
		}
	}

	slog.Info("Tables created successfully")
	return nil
}
