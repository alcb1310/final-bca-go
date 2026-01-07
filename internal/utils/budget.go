package utils

import (
	"database/sql"

	"github.com/alcb1310/final-bca-go/internal/types"
	"github.com/google/uuid"
)

func SaveBudget(b *types.CreateBudget, tx *sql.Tx) error {
	if b == nil || b.BudgetItemId == uuid.Nil {
		return nil
	}

	budget := types.SaveBudget{
		ProjectId:         b.ProjectId,
		BudgetItemId:      b.BudgetItemId,
		InitialQuantity:   sql.NullFloat64{Float64: b.Quantity, Valid: true},
		InitialCost:       sql.NullFloat64{Float64: b.Cost, Valid: true},
		InitialTotal:      b.Quantity * b.Cost,
		SpentQuantity:     sql.NullFloat64{Float64: 0, Valid: true},
		SpentTotal:        0,
		RemainingQuantity: sql.NullFloat64{Float64: b.Quantity, Valid: true},
		RemainingCost:     sql.NullFloat64{Float64: b.Cost, Valid: true},
		RemainingTotal:    b.Quantity * b.Cost,
		UpdatedBudget:     b.Quantity * b.Cost,
	}

	query := "select accumulate, parent_id from budget_item where id = $1"
	var accumulate bool
	var parentId uuid.UUID

	err := tx.QueryRow(query, budget.BudgetItemId).Scan(&accumulate, &parentId)
	if err != nil {
		return err
	}

	if !accumulate {
		query = `
			insert into budget
			(project_id, budget_item_id, initial_quantity, initial_cost,
			initial_total, spent_quantity, spent_total, remaining_quantity,
			remaining_cost, remaining_total, updated_budget)
			values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		`

		_, err = tx.Exec(query, budget.ProjectId, budget.BudgetItemId, budget.InitialQuantity,
			budget.InitialCost, budget.InitialTotal, budget.SpentQuantity,
			budget.SpentTotal, budget.RemainingQuantity, budget.RemainingCost,
			budget.RemainingTotal, budget.UpdatedBudget)
		if err != nil {
			return err
		}
	} else {
		query = `
			insert into budget 
			(project_id, budget_item_id, initial_total, spent_total, remaining_total, updated_budget)
			values ($1, $2, $3, $4, $5, $6) on conflict (project_id, budget_item_id)
			do update set initial_total = budget.initial_total + $3,
			spent_total = budget.spent_total + $4, remaining_total = budget.remaining_total + $5,
			updated_budget = budget.updated_budget + $6
		`

		_, err = tx.Exec(query, budget.ProjectId, budget.BudgetItemId, budget.InitialTotal,
			budget.SpentTotal, budget.RemainingTotal, budget.UpdatedBudget)
		if err != nil {
			return err
		}
	}

	b.BudgetItemId = parentId

	return SaveBudget(b, tx)
}
