package tests

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestApiInvoicesTest(t *testing.T) {
	ctx := context.Background()
	testUrl := "/api/v2/invoices"
	t.Log(testUrl)

	pgContainer, err := postgres.Run(ctx,
		"postgres:18-alpine",
		postgres.WithOrderedInitScripts(
			filepath.Join("..", "schema", "tables.sql"),
			filepath.Join("scripts", "seed_projects.sql"),
		),
		postgres.WithDatabase("testbca"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(15*time.Second)),
	)

	assert.NoError(t, err)
	if err != nil {
		if pgContainer != nil {
			err := pgContainer.Terminate(ctx)
			assert.NoError(t, err)
		}
		return
	}

	t.Cleanup(func() {
		err := pgContainer.Terminate(ctx)
		assert.NoError(t, err)
	})

	s, err := createServer(t, ctx, pgContainer)
	assert.NoError(t, err)
	if err != nil {
		return
	}
	s.GenerateRoutes()
}
