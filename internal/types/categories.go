package types

import "github.com/google/uuid"

type Category struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
