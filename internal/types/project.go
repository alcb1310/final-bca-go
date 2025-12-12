package types

import "github.com/google/uuid"

type Project struct {
	Id        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	IsActive  bool      `json:"is_active"`
	GrossArea float64   `json:"gross_area"`
	NetArea   float64   `json:"net_area"`
}
