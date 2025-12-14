package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/alcb1310/final-bca-go/internal/database"
	"github.com/alcb1310/final-bca-go/internal/router"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func createServer(t *testing.T, ctx context.Context, pgContainer *postgres.PostgresContainer) (*router.Router, error) {
	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	assert.NoError(t, err)
	db, _ := database.New(connStr)
	assert.NotNil(t, db)
	if db == nil {
		return nil, fmt.Errorf("db is nil")
	}

	s := router.NewRouter(db)
	return s, err
}
