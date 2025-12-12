package database

import (
	"log/slog"

	"github.com/alcb1310/final-bca-go/internal/types"
)

func (s *service) CreateProject(p types.Project) error {
	sql := "insert into project (name, is_active, gross_area, net_area) values ($1, $2, $3, $4)"
	if _, err := s.db.Exec(sql, p.Name, p.IsActive, p.GrossArea, p.NetArea); err != nil {
		slog.Error("Error creating project", "err", err)
	}
	return nil
}
