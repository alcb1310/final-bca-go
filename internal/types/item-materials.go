package types

import "github.com/google/uuid"

type ItemMaterialsResponse struct {
	ItemId       uuid.UUID `json:"item_id"`
	ItemCode     string    `json:"item_code"`
	ItemName     string    `json:"item_name"`
	ItemUnit     string    `json:"item_unit"`
	MaterialId   uuid.UUID `json:"material_id"`
	MaterialCode string    `json:"material_code"`
	MaterialName string    `json:"material_name"`
	MaterialUnit string    `json:"material_unit"`
	Quantity     float64   `json:"quantity"`
}
