package database

import (
	"log/slog"

	"github.com/alcb1310/final-bca-go/internal/types"
	"github.com/google/uuid"
)

func (s *service) CreateProject(p types.Project) error {
	sql := "insert into project (name, is_active, gross_area, net_area) values ($1, $2, $3, $4)"
	if _, err := s.db.Exec(sql, p.Name, p.IsActive, p.GrossArea, p.NetArea); err != nil {
		slog.Error("CreateProject: Error creating project", "err", err)
		return err
	}
	return nil
}

func (s *service) GetProjects() ([]types.Project, error) {
	var projects []types.Project
	sql := "select id, name, is_active, gross_area, net_area from project"
	rows, err := s.db.Query(sql)
	if err != nil {
		slog.Error("Error getting projects", "err", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var p types.Project
		if err := rows.Scan(&p.Id, &p.Name, &p.IsActive, &p.GrossArea, &p.NetArea); err != nil {
			slog.Error("GetProjects: Error scanning project", "err", err)
			return nil, err
		}
		projects = append(projects, p)
	}

	return projects, nil
}

func (s *service) GetProject(id uuid.UUID) (types.Project, error) {
	var err error
	project := types.Project{}

	sql := "select id, name, is_active, gross_area, net_area from project where id = $1"
	err = s.db.QueryRow(sql, id).Scan(&project.Id, &project.Name, &project.IsActive, &project.GrossArea, &project.NetArea)

	if err != nil {
		slog.Error("GetProject: Error scanning project", "err", err)
	}

	return project, err
}
