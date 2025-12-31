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
	GetCategory(id uuid.UUID) (types.Category, error)
	UpdateCategory(cat types.Category) error

	// file: materials.go
	CreateMaterial(mat types.Materials) error
	GetMaterials() ([]types.Materials, error)
	GetMaterial(id uuid.UUID) (types.Materials, error)
	UpdateMaterial(mat types.Materials) error
	DeleteMaterial(id uuid.UUID) error

	// file: items.go
	CreateItem(item *types.Items) error
	GetItems() ([]types.Items, error)
	GetItem(id uuid.UUID) (types.Items, error)
	UpdateItem(item types.Items) error

	// file item-materials.go
	GetItemMaterials(rubroId uuid.UUID) ([]types.ItemMaterialsResponse, error)
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
