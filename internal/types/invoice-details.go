package types

import "github.com/google/uuid"

type InvoiceDetailsResponse struct {
	Invoice    InvoiceResponse `json:"invoice"`
	Project    Project         `json:"project"`
	Supplier   Supplier        `json:"supplier"`
	BudgetItem BudgetItem      `json:"budget_item"`
	Quantity   float64         `json:"quantity"`
	Cost       float64         `json:"cost"`
	Total      float64         `json:"total"`
}

type InvoiceDetailsCreate struct {
	InvoiceId    uuid.UUID `json:"invoice_id"`
	BudgetItemId uuid.UUID `json:"budget_item_id"`
	Quantity     float64   `json:"quantity"`
	Cost         float64   `json:"cost"`
}
