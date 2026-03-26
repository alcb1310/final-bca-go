package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/alcb1310/final-bca-go/internal/types"
	"github.com/alcb1310/final-bca-go/internal/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
)

func (s *service) GenerateClosure(projectId uuid.UUID, date time.Time) error {
	var customError *utils.BcaError = nil

	tx, err := s.db.Begin()
	if err != nil {
		slog.Error("GenerateClosure: Error creating transaction", "err", err)
		customError = &utils.BcaError{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("Error creating transaction: %v", err),
		}

		return customError
	}
	defer func() {
		_ = tx.Rollback()
	}()

	query := "update project set last_closure = $1 where id = $2"
	_, err = tx.Exec(query, date, projectId)
	if err != nil {
		if err == sql.ErrNoRows {
			customError = &utils.BcaError{
				Code:    http.StatusNotFound,
				Message: "Project no encontrado",
			}
			return customError
		}

		slog.Error("GenerateClosure: Error updating project", "err", err)
		customError = &utils.BcaError{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("Error updating project: %v", err),
		}
		return customError
	}

	query = "select project_id, budget_item_id, initial_quantity, initial_cost, initial_total, spent_quantity, spent_total, remaining_quantity, remaining_cost, remaining_total, updated_budget from budget where project_id = $1"
	rows, err := tx.Query(query, projectId)
	if err != nil {
		if err == sql.ErrNoRows {
			customError = &utils.BcaError{
				Code:    http.StatusNotFound,
				Message: "El proyecto no tiene presupuesto",
			}
			return customError
		}

		slog.Error("GenerateClosure: Error getting budget", "err", err)
		customError = &utils.BcaError{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("Error getting budget: %v", err),
		}
		return customError
	}
	defer rows.Close()
	budgets := []types.SaveBudget{}

	for rows.Next() {
		b := types.SaveBudget{}

		if err := rows.Scan(&b.ProjectId, &b.BudgetItemId, &b.InitialQuantity, &b.InitialCost, &b.InitialTotal, &b.SpentQuantity, &b.SpentTotal, &b.RemainingQuantity, &b.RemainingCost, &b.RemainingTotal, &b.UpdatedBudget); err != nil {
			slog.Error("GenerateClosure: Error scanning budget", "err", err)
			customError = &utils.BcaError{
				Code:    http.StatusInternalServerError,
				Message: fmt.Sprintf("Error scanning budget: %v", err),
			}
			return customError
		}

		budgets = append(budgets, b)
	}

	for _, b := range budgets {
		sql := `
			insert into historic (
				project_id, budget_item_id, historic_date, initial_quantity, initial_cost,
				initial_total, spent_quantity, spent_total, remaining_quantity, remaining_cost,
				remaining_total, updated_budget)
			values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		`
		if _, err := tx.Exec(sql, b.ProjectId, b.BudgetItemId, date, b.InitialQuantity, b.InitialCost,
			b.InitialTotal, b.SpentQuantity, b.SpentTotal, b.RemainingQuantity, b.RemainingCost,
			b.RemainingTotal, b.UpdatedBudget); err != nil {
			var e *pgconn.PgError
			if errors.As(err, &e) {
				slog.Error("GenerateClosure: Error inserting historic", "err", e, "budget", b)
				customError = &utils.BcaError{
					Message: e.Message,
				}

				switch e.Code {
				case "23505":
					customError.Code = http.StatusConflict
				}

				return customError
			}

			slog.Error("GenerateClosure: Error inserting historic", "err", err, "budget", b)
			customError = &utils.BcaError{
				Code:    http.StatusInternalServerError,
				Message: fmt.Sprintf("Error inserting historic: %v", err),
			}
			return customError
		}
	}

	sql := "update invoice set is_balanced = true where project_id = $1 and extract(year from invoice_date) = $2 and extract(month from invoice_date) = $3"
	_, err = tx.Exec(sql, projectId, date.Year(), date.Month())
	if err != nil {
		slog.Error("GenerateClosure: Error updating invoices", "err", err)
		customError = &utils.BcaError{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("Error getting invoice: %v", err),
		}
		return customError
	}

	_ = tx.Commit()
	return customError
}
