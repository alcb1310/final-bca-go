package types

import (
	"database/sql"
)

type Budget struct {
	BudgetItem        BudgetItem      `json:"budget_item"`
	Project           Project         `json:"project"`
	InitialQuantity   sql.NullFloat64 `json:"initial_quantity"`
	InitialCost       sql.NullFloat64 `json:"initial_cost"`
	InitialTotal      float64         `json:"initial_total"`
	SpentQuantity     sql.NullFloat64 `json:"spent_quantity"`
	SpentTotal        float64         `json:"spent_total"`
	RemainingQuantity sql.NullFloat64 `json:"remaining_quantity"`
	RemainingCost     sql.NullFloat64 `json:"remaining_cost"`
	RemainingTotal    float64         `json:"remaining_total"`
	UpdatedBudget     float64         `json:"updated_budget"`
}
