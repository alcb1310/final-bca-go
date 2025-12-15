package types

import (
	"database/sql"

	"github.com/google/uuid"
)

type Supplier struct {
	Id           uuid.UUID      `json:"id"`
	Name         string         `json:"name"`
	SupplierId   string         `json:"supplier_id"`
	ContactName  sql.NullString `json:"contact_name"`
	ContactEmail sql.NullString `json:"contact_email"`
	ContactPhone sql.NullString `json:"contact_phone"`
}
