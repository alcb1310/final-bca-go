package database

import (
	"github.com/alcb1310/final-bca-go/internal/types"
	"github.com/google/uuid"
)

func (s *service) ListInvoiceDetails(id uuid.UUID) ([]types.InvoiceDetailsResponse, error) {
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
			supplier_number
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
