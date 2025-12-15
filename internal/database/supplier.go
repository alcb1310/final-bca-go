package database

import (
	"log/slog"

	"github.com/alcb1310/final-bca-go/internal/types"
)

func (s *service) GetSuppliers() ([]types.Supplier, error) {
	suppliers := []types.Supplier{}
	sql := "select id, name, supplier_id, contact_name, contac_phone, contact_email from supplier"
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
