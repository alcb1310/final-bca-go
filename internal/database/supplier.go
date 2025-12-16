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

func (s *service) CreateSupplier(su types.Supplier) error {
	var err error
	sql := "insert into supplier (name, supplier_id, contact_name, contact_phone, contact_email) values ($1, $2, $3, $4, $5)"
	if _, err = s.db.Exec(sql, su.Name, su.SupplierId, su.ContactName, su.ContactPhone, su.ContactEmail); err != nil {
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
