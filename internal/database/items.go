package database

import "github.com/alcb1310/final-bca-go/internal/types"

func (s *service) GetItems() ([]types.Items, error) {
	items := []types.Items{}

	sql := "select id, code, name, unit from item order by name"
	rows, err := s.db.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item types.Items
		if err := rows.Scan(&item.Id, &item.Code, &item.Name, &item.Unit); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}
