package database

import (
	"log/slog"

	"github.com/alcb1310/final-bca-go/internal/types"
	"github.com/google/uuid"
)

func (s *service) GetInvoiceDetails(id uuid.UUID) ([]types.InvoiceDetailsResponse, error) {
	invoiceDetails := []types.InvoiceDetailsResponse{}

	query := `
		SELECT
			invoice_id,
			invoice_number,
			invoice_total,
			invoice_date,
			project_id,
			project_name,
			supplier_id,
			supplier_number,
			supplier_name,
			budget_item_id,
			budget_item_code,
			budget_item_name,
			budget_item_level,
			quantity,
			cost,
			total
		FROM vw_invoice_details
		WHERE invoice_id = $1
	`
	rows, err := s.db.Query(query, id)
	if err != nil {
		return invoiceDetails, err
	}
	defer rows.Close()

	for rows.Next() {
		id := types.InvoiceDetailsResponse{}
		if err := rows.Scan(
			&id.Invoice.Id,
			&id.Invoice.InvoiceNumber,
			&id.Invoice.InvoiceTotal,
			&id.Invoice.InvoiceDate,
			&id.Project.Id,
			&id.Project.Name,
			&id.Supplier.Id,
			&id.Supplier.SupplierId,
			&id.Supplier.Name,
			&id.BudgetItem.Id,
			&id.BudgetItem.Code,
			&id.BudgetItem.Name,
			&id.BudgetItem.Level,
			&id.Quantity,
			&id.Cost,
			&id.Total,
		); err != nil {
			return invoiceDetails, err
		}

		invoiceDetails = append(invoiceDetails, id)
	}

	return invoiceDetails, nil
}

func (s *service) GetInvoiceDetail(invoiceId, budgetItemId uuid.UUID) (types.InvoiceDetailsResponse, error) {
	id := types.InvoiceDetailsResponse{}
	query := `
		SELECT
			invoice_id,
			invoice_number,
			invoice_total,
			invoice_date,
			project_id,
			project_name,
			supplier_id,
			supplier_number,
			supplier_name,
			budget_item_id,
			budget_item_code,
			budget_item_name,
			budget_item_level,
			quantity,
			cost,
			total
		FROM vw_invoice_details
		WHERE invoice_id = $1 and budget_item_id = $2
	`
	err := s.db.QueryRow(query, invoiceId, budgetItemId).Scan(
		&id.Invoice.Id,
		&id.Invoice.InvoiceNumber,
		&id.Invoice.InvoiceTotal,
		&id.Invoice.InvoiceDate,
		&id.Project.Id,
		&id.Project.Name,
		&id.Supplier.Id,
		&id.Supplier.SupplierId,
		&id.Supplier.Name,
		&id.BudgetItem.Id,
		&id.BudgetItem.Code,
		&id.BudgetItem.Name,
		&id.BudgetItem.Level,
		&id.Quantity,
		&id.Cost,
		&id.Total,
	)

	return id, err
}

func (s *service) CreateInvoiceDetail(detail types.InvoiceDetailsCreate) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()
	total := detail.Quantity * detail.Cost

	query := "insert into invoice_details (invoice_id, budget_item_id, quantity, cost, total) values ($1, $2, $3, $4, $5)"
	_, err = tx.Exec(query, detail.InvoiceId, detail.BudgetItemId, detail.Quantity, detail.Cost, total)
	if err != nil {
		slog.Error("Error creating invoice detail insert", "err", err)
		_ = tx.Rollback()
		return err
	}
	query = "update invoice set invoice_total = invoice_total + $1 where id = $2"
	_, err = tx.Exec(query, total, detail.InvoiceId)
	if err != nil {
		slog.Error("Error creating invoice detail update", "err", err)
		_ = tx.Rollback()
		return err
	}

	var projectId uuid.UUID
	query = "select project_id from invoice where id = $1"
	if err := tx.QueryRow(query, detail.InvoiceId).Scan(&projectId); err != nil {
		slog.Error("Error getting project id", "err", err)
		_ = tx.Rollback()
		return err
	}

	query = "select spent_quantity, spent_total, remaining_quantity, remaining_cost, remaining_total, updated_budget from budget where project_id = $1 and budget_item_id = $2"
	var spentQuantity, spentTotal, remainingQuantity, remainingCost, remainingTotal, updatedBudget float64
	if err := tx.QueryRow(query, projectId, detail.BudgetItemId).Scan(&spentQuantity, &spentTotal, &remainingQuantity, &remainingCost, &remainingTotal, &updatedBudget); err != nil {
		slog.Error("Error getting budget", "err", err)
		_ = tx.Rollback()
		return err
	}

	newToSpendTotal := (remainingQuantity - detail.Quantity) * detail.Cost
	newSpentTotal := spentTotal + total
	newUpdatedBudget := newSpentTotal + newToSpendTotal

	query = `
		update budget set 
		spent_quantity = spent_quantity + $1,
		spent_total = spent_total + $3,
		remaining_quantity = remaining_quantity - $1,
		remaining_cost = $2,
		remaining_total = $4,
		updated_budget = $5
		where project_id = $6 and budget_item_id = $7
	`

	if _, err = tx.Exec(query, detail.Quantity, detail.Cost, total, newToSpendTotal, newUpdatedBudget, projectId, detail.BudgetItemId); err != nil {
		slog.Error("Error updating budget", "err", err)
		_ = tx.Rollback()
		return err
	}

	updatedDiff := newUpdatedBudget - updatedBudget
	remainingDiff := newToSpendTotal - remainingTotal

	parentId := &detail.BudgetItemId
	for {
		query = "select parent_id from budget_item where id = $1"
		if err := tx.QueryRow(query, parentId).Scan(&parentId); err != nil {
			slog.Error("Error getting parent id", "err", err)
			_ = tx.Rollback()
			return err
		}

		if parentId == nil || parentId == &uuid.Nil {
			break
		}

		query = `
			update budget set
			spent_total = spent_total + $1,
			remaining_total = remaining_total + $2,
			updated_budget = updated_budget + $3
			where project_id = $4 and budget_item_id = $5
		`
		if _, err = tx.Exec(query, total, remainingDiff, updatedDiff, projectId, *parentId); err != nil {
			slog.Error("Error updating parent budget", "err", err)
			_ = tx.Rollback()
			return err
		}
	}

	slog.Info("Invoice detail created successfully")
	_ = tx.Commit()

	return nil
}

func (s *service) DeleteInvoiceDetail(invoiceId uuid.UUID, budgetItemId uuid.UUID) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	var quantity, cost, total float64
	query := "select quantity, cost, total from invoice_details where invoice_id = $1 and budget_item_id = $2"
	if err := tx.QueryRow(query, invoiceId, budgetItemId).Scan(&quantity, &cost, &total); err != nil {
		slog.Error("Error getting invoice detail", "err", err)
		return err
	}

	query = "delete from invoice_details where invoice_id = $1 and budget_item_id = $2"
	if _, err := tx.Exec(query, invoiceId, budgetItemId); err != nil {
		slog.Error("Error deleting invoice detail", "err", err)
		return err
	}

	query = "update invoice set invoice_total = invoice_total - $1 where id = $2"
	if _, err := tx.Exec(query, total, invoiceId); err != nil {
		slog.Error("Error updating invoice total", "err", err)
		return err
	}

	var projectId uuid.UUID
	query = "select project_id from invoice where id = $1"
	if err := tx.QueryRow(query, invoiceId).Scan(&projectId); err != nil {
		slog.Error("Error getting project id", "err", err)
		return err
	}

	query = `
		update budget set
			spent_quantity = spent_quantity - $1,
			spent_total = spent_total - $2,
			remaining_quantity = remaining_quantity + $1,
			remaining_cost = $3,
			remaining_total = remaining_total + $2
		where project_id = $4 and budget_item_id = $5
	`
	if _, err := tx.Exec(query, quantity, total, cost, projectId, budgetItemId); err != nil {
		slog.Error("Error updating budget", "err", err)
		return err
	}

	parentId := &budgetItemId
	for {
		query = "select parent_id from budget_item where id = $1"
		if err := tx.QueryRow(query, parentId).Scan(&parentId); err != nil {
			slog.Error("Error getting parent id", "err", err)
			return err
		}
		if parentId == nil || parentId == &uuid.Nil {
			break
		}

		query = `
			update budget set
				spent_total = spent_total - $1,
				remaining_total = remaining_total + $1
			where project_id = $2 and budget_item_id = $3
		`

		if _, err := tx.Exec(query, total, projectId, parentId); err != nil {
			slog.Error("Error updating parent budget", "err", err)
			return err
		}
	}

	_ = tx.Commit()
	return nil
}
