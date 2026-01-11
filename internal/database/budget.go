package database

import (
	"database/sql"
	"log/slog"

	"github.com/alcb1310/final-bca-go/internal/types"
	"github.com/alcb1310/final-bca-go/internal/utils"
	"github.com/google/uuid"
)

func (s *service) GetBudgets() ([]types.Budget, error) {
	budgets := []types.Budget{}

	sql := `
		select
			budget_item_id, budget_item_code, budget_item_name, budget_item_level, budget_item_accumulate,
			project_id, project_name,
			initial_quantity, initial_cost, initial_total,
			spent_quantity, spent_total,
			remaining_quantity, remaining_cost, remaining_total,
			updated_budget
		from vw_budget
		order by project_name, budget_item_code
	`

	rows, err := s.db.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		b := types.Budget{}

		if err := rows.Scan(
			&b.BudgetItem.Id, &b.BudgetItem.Code, &b.BudgetItem.Name, &b.BudgetItem.Level, &b.BudgetItem.Accumulate,
			&b.Project.Id, &b.Project.Name,
			&b.InitialQuantity, &b.InitialCost, &b.InitialTotal,
			&b.SpentQuantity, &b.SpentTotal,
			&b.RemainingQuantity, &b.RemainingCost, &b.RemainingTotal,
			&b.UpdatedBudget,
		); err != nil {
			return nil, err
		}

		budgets = append(budgets, b)
	}

	return budgets, nil
}

func (s *service) GetBudget(projectId uuid.UUID, budgetItemId uuid.UUID) (types.SaveBudget, error) {
	var b types.SaveBudget
	query := `
		select project_id, budget_item_id, initial_quantity, initial_cost, initial_total,
		spent_quantity, spent_total, remaining_quantity, remaining_cost, remaining_total,
		updated_budget
		from budget
		where project_id = $1 and budget_item_id = $2
	`

	err := s.db.QueryRow(query, projectId, budgetItemId).Scan(
		&b.ProjectId, &b.BudgetItemId, &b.InitialQuantity, &b.InitialCost, &b.InitialTotal,
		&b.SpentQuantity, &b.SpentTotal, &b.RemainingQuantity, &b.RemainingCost, &b.RemainingTotal,
		&b.UpdatedBudget,
	)

	return b, err
}

func (s *service) CreateBudget(b types.CreateBudget) error {
	tx, err := s.db.Begin()
	if err != nil {
		slog.Error("Error creating transaction", "err", err)
		return err
	}
	defer func() { _ = tx.Commit() }()

	if err := utils.SaveBudget(&b, tx); err != nil {
		_ = tx.Rollback()
		slog.Error("Error saving budget", "err", err)
		return err
	}

	return nil
}

func (s *service) UpdateBudget(b types.CreateBudget, oldBudget types.SaveBudget) error {
	total := b.Quantity * b.Cost
	diff := total - oldBudget.RemainingTotal

	q := sql.NullFloat64{Float64: b.Quantity, Valid: true}
	c := sql.NullFloat64{Float64: b.Cost, Valid: true}

	toUpdate := types.SaveBudget{
		ProjectId:         oldBudget.ProjectId,
		BudgetItemId:      oldBudget.BudgetItemId,
		InitialQuantity:   oldBudget.InitialQuantity,
		InitialCost:       oldBudget.InitialCost,
		InitialTotal:      oldBudget.InitialTotal,
		SpentQuantity:     oldBudget.SpentQuantity,
		SpentTotal:        oldBudget.SpentTotal,
		RemainingQuantity: q,
		RemainingCost:     c,
		RemainingTotal:    total,
		UpdatedBudget:     total + oldBudget.SpentTotal,
	}

	tx, err := s.db.Begin()
	if err != nil {
		slog.Error("Error creating transaction", "err", err)
		return err
	}
	defer func() { _ = tx.Commit() }()

	if err := s.runUpdateBudget(&toUpdate, tx, diff); err != nil {
		_ = tx.Rollback()
		slog.Error("Error updating budget", "err", err)
		return err
	}

	return nil
}

func (s *service) runUpdateBudget(budgets *types.SaveBudget, tx *sql.Tx, diff float64) error {
	nullUUID := uuid.UUID{}
	bi, err := s.GetBudgetItem(budgets.BudgetItemId)
	if err != nil {
		return err
	}
	if err == sql.ErrNoRows {
		return nil
	}

	if bi.Accumulate {
		query := `
			UPDATE budget
			SET remaining_total = budget.remaining_total + $1,
			updated_budget = budget.updated_budget + $1
			WHERE project_id = $2 AND budget_item_id = $3
		`
		_, err = tx.Exec(query, diff, budgets.ProjectId, budgets.BudgetItemId)
	} else {
		query := `
			UPDATE budget
			SET remaining_quantity = $1, remaining_cost = $2, remaining_total = $3,
			updated_budget = $4
			WHERE project_id = $5 AND budget_item_id = $6
		`
		_, err = tx.Exec(query, budgets.RemainingQuantity, budgets.RemainingCost, budgets.RemainingTotal,
			budgets.UpdatedBudget, budgets.ProjectId, budgets.BudgetItemId)
	}

	if err != nil {
		return err
	}
	if bi.ParentId == nullUUID {
		return nil
	}
	budgets.BudgetItemId = bi.ParentId

	return s.runUpdateBudget(budgets, tx, diff)
}
