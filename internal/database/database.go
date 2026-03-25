package database

import (
	"database/sql"
	"fmt"
	"os"
	"time"

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
	CreateItemMaterial(im types.ItemMaterialCreate) error
	GetItemMaterials(rubroId uuid.UUID) ([]types.ItemMaterialsResponse, error)
	UpdateItemMaterial(im types.ItemMaterialCreate) error
	DeleteItemMaterial(rubroId, materialId uuid.UUID) error

	// file budget.go
	CreateBudget(budget types.CreateBudget) error
	GetBudgets() ([]types.Budget, error)
	GetBudget(projectId uuid.UUID, Id uuid.UUID) (types.SaveBudget, error)
	UpdateBudget(budget types.CreateBudget, oldBudget types.SaveBudget) error

	// file invoice.go
	CreateInvoice(inv *types.InvoiceCreate) error
	GetInvoices() ([]types.InvoiceResponse, error)
	GetInvoice(id uuid.UUID) (types.InvoiceResponse, error)
	UpdateInvoice(inv types.InvoiceUpdate) error
	DeleteInvoice(id uuid.UUID) error

	// file invoice-detail.go
	CreateInvoiceDetail(detail types.InvoiceDetailsCreate) error
	GetInvoiceDetails(id uuid.UUID) ([]types.InvoiceDetailsResponse, error)
	GetInvoiceDetail(invoiceId, budgetItemId uuid.UUID) (types.InvoiceDetailsResponse, error)
	DeleteInvoiceDetail(invoiceId, budgetItemId uuid.UUID) error

	// file closure.go
	GenerateClosure(projectId uuid.UUID, data time.Time) error
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
