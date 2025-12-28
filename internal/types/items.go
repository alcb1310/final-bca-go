package types

import "github.com/google/uuid"

type Items struct {
	Id   uuid.UUID `json:"id"`
	Code string    `json:"code"`
	Name string    `json:"name"`
	Unit string    `json:"unit"`
}
