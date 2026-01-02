package database

import (
	"github.com/alcb1310/final-bca-go/internal/types"
	"github.com/google/uuid"
)

func (s *service) GetItemMaterials(rubroId uuid.UUID) ([]types.ItemMaterialsResponse, error) {
	im := []types.ItemMaterialsResponse{}

	sql := "select material_id, material_code, material_name, material_unit, item_id, item_code, item_name, item_unit from vw_item_materials where item_id = $1"

	rows, err := s.db.Query(sql, rubroId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var imr types.ItemMaterialsResponse
		if err := rows.Scan(&imr.MaterialId, &imr.MaterialCode, &imr.MaterialName, &imr.MaterialUnit, &imr.ItemId, &imr.ItemCode, &imr.ItemName, &imr.ItemUnit); err != nil {
			return nil, err
		}
		im = append(im, imr)
	}

	return im, nil
}

func (s *service) CreateItemMaterial(im types.ItemMaterialCreate) error {
	sql := "insert into item_material (material_id, item_id) values ($1, $2)"
	_, err := s.db.Exec(sql, im.MaterialId, im.ItemId)
	return err
}

func (s *service) UpdateItemMaterial(im types.ItemMaterialCreate) error {
	sql := "update item_material set quantity = $1 where material_id = $2 and item_id = $3"
	_, err := s.db.Exec(sql, im.Quantity, im.MaterialId, im.ItemId)
	return err
}
