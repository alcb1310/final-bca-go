package database

import (
	"github.com/alcb1310/final-bca-go/internal/types"
	"github.com/google/uuid"
)

func (s *service) GetMaterials() ([]types.Materials, error) {
	materials := []types.Materials{}

	sql := "select id, code, name, unit, category_id, category_name from vw_materials"
	rows, err := s.db.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var material types.Materials
		if err := rows.Scan(&material.Id, &material.Code, &material.Name, &material.Unit, &material.Category.Id, &material.Category.Name); err != nil {
			return nil, err
		}
		materials = append(materials, material)
	}

	return materials, nil
}

func (s *service) CreateMaterial(mat types.Materials) error {
	sql := "insert into material (code, name, unit, category_id) values ($1, $2, $3, $4)"
	_, err := s.db.Exec(sql, mat.Code, mat.Name, mat.Unit, mat.Category.Id)
	return err
}

func (s *service) GetMaterial(id uuid.UUID) (types.Materials, error) {
	material := types.Materials{}

	sql := "select id, code, name, unit, category_id, category_name from vw_materials where id = $1"
	err := s.db.QueryRow(sql, id).Scan(&material.Id, &material.Code, &material.Name, &material.Unit, &material.Category.Id, &material.Category.Name)
	return material, err
}

func (s *service) UpdateMaterial(mat types.Materials) error {
	sql := "update material set code = $1, name = $2, unit = $3, category_id = $4 where id = $5"
	_, err := s.db.Exec(sql, mat.Code, mat.Name, mat.Unit, mat.Category.Id, mat.Id)
	return err
}
