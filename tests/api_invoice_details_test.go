package tests

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestApiInvoiceDetails(t *testing.T) {
	ctx := context.Background()
	pgContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithOrderedInitScripts(
			filepath.Join("..", "schema", "tables.sql"),
			filepath.Join("scripts", "seed_projects.sql"),
			filepath.Join("scripts", "seed_suppliers.sql"),
			filepath.Join("scripts", "seed_budget-items.sql"),
			filepath.Join("scripts", "seed_budget.sql"),
			filepath.Join("scripts", "seed_invoices.sql"),
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
		slog.Error("TestApiInvoiceDetails, failed to run pgContainer", "error", err)
		panic(err)
	}

	t.Cleanup(func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			t.Fatalf("TestApiInvoiceDetails, failed to terminate pgContainer: %v", err)
		}
	})

	s, err := createServer(t, ctx, pgContainer)
	assert.NoError(t, err)
	if err != nil {
		slog.Error("TestApiInvoiceDetails, failed to create server", "error", err)
		panic(err)
	}
	s.GenerateRoutes()
	invoiceId := uuid.MustParse("c3be2956-1c3c-46f7-af14-d28420116f14")
	testURL := fmt.Sprintf("/api/v2/invoices/%s/details", invoiceId)

	t.Run("should have no invoice details", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, testURL, nil)
		assert.NoError(t, err)
		resp := httptest.NewRecorder()
		s.Router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Equal(t, "[]", strings.TrimSpace(resp.Body.String()))
	})
}
