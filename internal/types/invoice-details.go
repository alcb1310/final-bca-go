package types

import "github.com/google/uuid"

type InvoiceDetails struct {
	InvoiceId      uuid.UUID `json:"invoice_id"`
	InvoiceNumber  string    `json:"invoice_number"`
	ProjectId      uuid.UUID `json:"project_id"`
	ProjectName    string    `json:"project_name"`
	SupplierId     uuid.UUID `json:"supplier_id"`
	SupplierName   string    `json:"supplier_name"`
	BudgetItemId   uuid.UUID `json:"budget_item_id"`
	BudgetItemCode string    `json:"budget_item_code"`
	BudgetItemName string    `json:"budget_item_name"`
	Quantity       float64   `json:"quantity"`
	Cost           float64   `json:"cost"`
	Total          float64   `json:"total"`
}
