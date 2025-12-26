package database

import "github.com/alcb1310/final-bca-go/internal/types"

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
