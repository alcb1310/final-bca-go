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

	// file: project.go
	CreateProject(p types.Project) error
	GetProjects() ([]types.Project, error)
	GetProject(id uuid.UUID) (types.Project, error)
	UpdateProject(p types.Project) error

	// file: supplier.go
	CreateSupplier(sup types.Supplier) error
	GetSuppliers() ([]types.Supplier, error)
	GetSupplier(id uuid.UUID) (types.Supplier, error)
	UpdateSupplier(sup types.Supplier) error

	// file: budget-item.go
	CreateBudgetItem(bi types.CreateBudgetItem) error
	GetBudgetItems() ([]types.BudgetItem, error)
	GetBudgetItemsByAccumulate(accum bool) ([]types.BudgetItem, error)
	GetBudgetItem(id uuid.UUID) (types.BudgetItem, error)
	UpdateBudgetItem(bi types.UpdateBudgetItem) error

	// file: categories.go
	CreateCategory(cat types.Category) error
	GetCategories() ([]types.Category, error)
}

type service struct {
	db *sql.DB
}

func New(connStr string) (Service, *sql.DB) {
	var db *sql.DB
	var err error
	if connStr == "" {
		fmt.Fprintf(os.Stderr, "New Router: DATABASE_URL is not set\n")
		return nil, nil
	}
	if db, err = sql.Open("pgx", connStr); err != nil {
		fmt.Fprintf(os.Stderr, "New Database: Unable to connect to database: %v\n", err)
		return nil, nil
	}

	s := &service{
		db: db,
	}
	return s, db
}
