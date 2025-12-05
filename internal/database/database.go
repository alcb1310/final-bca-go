package database

import (
)

type Service interface{}

type service struct {
}

func New() Service {
	var db *sql.DB
	s := &service{
		db: db,
	}
	return s
}
