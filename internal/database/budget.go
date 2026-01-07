package database

import "github.com/alcb1310/final-bca-go/internal/types"

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
