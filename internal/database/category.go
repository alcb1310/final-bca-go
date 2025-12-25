package database

import "github.com/alcb1310/final-bca-go/internal/types"

func (s *service) GetCategories() ([]types.Category, error) {
	categories := []types.Category{}

	sql := "select id, name from category order by name"
	rows, err := s.db.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var category types.Category
		if err := rows.Scan(&category.Id, &category.Name); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}

func (s *service) CreateCategory(cat types.Category) error {
	sql := "insert into category (name) values ($1)"
	_, err := s.db.Exec(sql, cat.Name)
	return err
}
