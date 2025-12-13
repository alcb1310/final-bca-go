package database

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/alcb1310/final-bca-go/internal/types"
	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Service interface {
	GetHealth() bool

	// file project.go
	CreateProject(p types.Project) error
	GetProjects() ([]types.Project, error)
	GetProject(id uuid.UUID) (types.Project, error)
}

type service struct {
	db *sql.DB
}

func New() Service {
	var db *sql.DB
	var err error
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		fmt.Fprintf(os.Stderr, "New Router: DATABASE_URL is not set\n")
		return nil
	}
	if db, err = sql.Open("pgx", connStr); err != nil {
		fmt.Fprintf(os.Stderr, "New Database: Unable to connect to database: %v\n", err)
		return nil
	}

	if err = createTables(db); err != nil {
		fmt.Fprintf(os.Stderr, "New Database: Unable to create tables: %v\n", err)
		return nil
	}

	s := &service{
		db: db,
	}
	return s
}
