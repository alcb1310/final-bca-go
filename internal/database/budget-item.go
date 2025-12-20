package database

import (
	"log/slog"

	"github.com/alcb1310/final-bca-go/internal/types"
)

func (s *service) GetBudgetItems() ([]types.BudgetItem, error) {
	bi := []types.BudgetItem{}

	sql := "select id, code, name, level, accumulate, parent_id, parent_code, parent_name from budget_items sort by code"
	rows, err := s.db.Query(sql)
	if err != nil {
		slog.Error("Error getting budget items", "err", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var b types.BudgetItem
		if err := rows.Scan(&b.Id, &b.Code, &b.Name, &b.Level, &b.Accumulate, &b.ParentId, &b.ParentCode, &b.ParentName); err != nil {
			slog.Error("GetBudgetItems: Error scanning budget item", "err", err)
			return nil, err
		}
		bi = append(bi, b)
	}

	return bi, nil
}
