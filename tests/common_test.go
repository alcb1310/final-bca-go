package tests

import (
	"context"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/alcb1310/final-bca-go/internal/database"
	"github.com/alcb1310/final-bca-go/internal/router"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func createServer(t *testing.T, ctx context.Context, pgContainer *postgres.PostgresContainer) (*router.Router, *httptest.Server, error) {
	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	assert.NoError(t, err)
	db, _ := database.New(connStr)
	assert.NotNil(t, db)
	if db == nil {
		return nil, nil, fmt.Errorf("db is nil")
	}

	server := &router.Router{
		DB: db,
	}

	s := httptest.NewServer(server.Router())

	return server, s, err
}
