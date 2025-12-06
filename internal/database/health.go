package database

import (
	"context"
	"log/slog"
	"time"
)

func (s *service) GetHealth() bool {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := s.db.PingContext(ctx)
	if err != nil {
		slog.Error("db down", "err", err)
		return false
	}
	return true
}
