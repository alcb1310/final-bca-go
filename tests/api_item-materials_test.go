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

func TestApiItemMaterialsTest(t *testing.T) {
	ctx := context.Background()
	testUrl := "/api/v2/items/2d257121-43e8-4b00-947d-b05fa54b36ac/materials"
	t.Log(testUrl)

	pgContainer, err := postgres.Run(ctx,
		"postgres:18-alpine",
		postgres.WithOrderedInitScripts(
			filepath.Join("..", "schema", "tables.sql"),
			filepath.Join("scripts", "seed_categories.sql"),
			filepath.Join("scripts", "seed_items.sql"),
			filepath.Join("scripts", "seed_materials.sql"),
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
