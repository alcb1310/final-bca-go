package database

import (
	"log/slog"

	"github.com/alcb1310/final-bca-go/internal/types"
)

func (s *service) GetBudgetItems() ([]types.BudgetItem, error) {
	bi := []types.BudgetItem{}

	sql := "select id, code, name, level, accumulate, parent_id, parent_code, parent_name from vw_budget_item order by code"
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

func (s *service) CreateBudgetItem(bi types.CreateBudgetItem) error {
	var level uint8 = 1
	var err error
	if bi.ParentId.Valid {
		sql := "select level from budget_item where id = $1"
		err := s.db.QueryRow(sql, bi.ParentId.UUID).Scan(&level)
		if err != nil {
			slog.Error("CreateBudgetItem: Error scanning budget item", "err", err)
			return err
		}
		level++
	}

	sql := "insert into budget_item(code, name, level, accumulate, parent_id) values ($1, $2, $3, $4, $5)"
	if bi.ParentId.Valid {
		_, err = s.db.Exec(sql, bi.Code, bi.Name, level, bi.Accumulate, bi.ParentId.UUID)
	} else {
		_, err = s.db.Exec(sql, bi.Code, bi.Name, level, bi.Accumulate, nil)
	}

	return err
}
