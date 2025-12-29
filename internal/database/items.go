package database

import (
	"github.com/alcb1310/final-bca-go/internal/types"
	"github.com/google/uuid"
)

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

func (s *service) CreateItem(item types.Items) error {
	sql := "insert into item (code, name, unit) values ($1, $2, $3)"
	_, err := s.db.Exec(sql, item.Code, item.Name, item.Unit)
	return err
}

func (s *service) GetItem(id uuid.UUID) (types.Items, error) {
	item := types.Items{}
	sql := "select id, code, name, unit from item where id = $1"
	err := s.db.QueryRow(sql, id).Scan(&item.Id, &item.Code, &item.Name, &item.Unit)
	return item, err
}

func (s *service) UpdateItem(item types.Items) error {
	sql := "update item set code = $1, name = $2, unit = $3 where id = $4"
	_, err := s.db.Exec(sql, item.Code, item.Name, item.Unit, item.Id)
	return err
}
