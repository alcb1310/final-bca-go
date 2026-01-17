package database

import (
	"github.com/alcb1310/final-bca-go/internal/types"
	"github.com/google/uuid"
)

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

func (s *service) GetInvoice(id uuid.UUID) (types.InvoiceResponse, error) {
	inv := types.InvoiceResponse{}
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
		WHERE id = $1
	`
	err := s.db.QueryRow(sql, id).Scan(
		&inv.Id,
		&inv.Supplier.Id,
		&inv.Supplier.SupplierId,
		&inv.Supplier.Name,
		&inv.Supplier.ContactName,
		&inv.Supplier.ContactEmail,
		&inv.Supplier.ContactPhone,
		&inv.Project.Id,
		&inv.Project.Name,
		&inv.Project.IsActive,
		&inv.InvoiceNumber,
		&inv.InvoiceDate,
		&inv.InvoiceTotal,
		&inv.IsBalanced,
	)
	return inv, err
}

func (s *service) CreateInvoice(inv types.InvoiceCreate) error {
	q := "INSERT INTO invoice (supplier_id, project_id, invoice_number, invoice_date) VALUES ($1, $2, $3, $4)"
	_, err := s.db.Exec(q, inv.SupplierId, inv.ProjectId, inv.InvoiceNumber, inv.InvoiceDate)
	return err
}

func (s *service) DeleteInvoice(id uuid.UUID) error {
	sql := "delete from invoice where id = $1"
	_, err := s.db.Exec(sql, id)
	return err
}

func (s *service) UpdateInvoice(inv types.InvoiceUpdate) error {
	sql := "update invoice set invoice_number = $1, invoice_date = $2 where id = $3"
	_, err := s.db.Exec(sql, inv.InvoiceNumber, inv.InvoiceDate, inv.Id)
	return err
}
