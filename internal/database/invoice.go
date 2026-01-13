package database

import "github.com/alcb1310/final-bca-go/internal/types"

func (s *service) GetInvoices() ([]types.InvoiceResponse, error) {
	invoices := []types.InvoiceResponse{}

	sql := `
		SELECT
			id,
			supplier_id,
			supplier_number,
			supplier_name,
			contact_name,
			contact_email,
			contact_phone,
			project_id,
			project_name,
			is_active,
			invoice_number,
			invoice_date,
			invoice_total,
			is_balanced
		FROM vw_invoice
	`
	rows, err := s.db.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var i types.InvoiceResponse
		if err := rows.Scan(
			&i.Id,
			&i.Supplier.Id,
			&i.Supplier.SupplierId,
			&i.Supplier.Name,
			&i.Supplier.ContactName,
			&i.Supplier.ContactEmail,
			&i.Supplier.ContactPhone,
			&i.Project.Id,
			&i.Project.Name,
			&i.Project.IsActive,
			&i.InvoiceNumber,
			&i.InvoiceDate,
			&i.InvoiceTotal,
			&i.IsBalanced,
		); err != nil {
			return nil, err
		}
		invoices = append(invoices, i)
	}

	return invoices, nil
}
