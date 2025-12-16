package database

import (
	"log/slog"

	"github.com/alcb1310/final-bca-go/internal/types"
	"github.com/google/uuid"
)

func (s *service) GetSuppliers() ([]types.Supplier, error) {
	suppliers := []types.Supplier{}
	sql := "select id, name, supplier_id, contact_name, contact_phone, contact_email from supplier"
	rows, err := s.db.Query(sql)
	if err != nil {
		slog.Error("GetAllSuppliers: error fetchingh the suppliers", "err", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var supplier types.Supplier
		if err := rows.Scan(&supplier.Id, &supplier.Name, &supplier.SupplierId, &supplier.ContactName, &supplier.ContactPhone, &supplier.ContactEmail); err != nil {
			slog.Error("GetAllSuppliers: error scanning the supplier", "err", err)
			return nil, err
		}
		suppliers = append(suppliers, supplier)
	}

	return suppliers, nil
}

func (s *service) CreateSupplier(sup types.Supplier) error {
	var err error
	sql := "insert into supplier (name, supplier_id, contact_name, contact_phone, contact_email) values ($1, $2, $3, $4, $5)"
	if _, err = s.db.Exec(sql, sup.Name, sup.SupplierId, sup.ContactName, sup.ContactPhone, sup.ContactEmail); err != nil {
		slog.Error("CreateSupplier: Error creating supplier", "err", err)
	}
	return err
}

func (s *service) GetSupplier(id uuid.UUID) (types.Supplier, error) {
	sup := types.Supplier{}

	sql := "select id, name, supplier_id, contact_name, contact_phone, contact_email from supplier where id = $1"
	err := s.db.QueryRow(sql, id).Scan(&sup.Id, &sup.Name, &sup.SupplierId, &sup.ContactName, &sup.ContactPhone, &sup.ContactEmail)

	return sup, err
}

func (s *service) UpdateSupplier(sup types.Supplier) error {
	var err error

	sql := "update supplier set name = $1, supplier_id = $2, contact_name = $3, contact_phone = $4, contact_email = $5 where id = $6"
	if _, err = s.db.Exec(sql, sup.Name, sup.SupplierId, sup.ContactName, sup.ContactPhone, sup.ContactEmail, sup.Id); err != nil {
		slog.Error("UpdateSupplier: Error updating supplier", "err", err)
	}

	return err
}
