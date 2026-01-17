package tests

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
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
			filepath.Join("scripts", "seed_suppliers.sql"),
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

	t.Run("should have no invoices", func(t *testing.T) {
		req, err := http.NewRequest("GET", testUrl, nil)
		assert.NoError(t, err)
		res := httptest.NewRecorder()
		s.Router.ServeHTTP(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
		var r []any
		err = json.Unmarshal(res.Body.Bytes(), &r)
		assert.NoError(t, err)
		assert.Equal(t, 0, len(r))
		assert.Equal(t, "[]", strings.TrimSpace(res.Body.String()))
	})

	t.Run("should be able to create an invoice", func(t *testing.T) {
		form := map[string]any{
			"project_id":     "1c6020db-39a0-451d-89ee-fdd20d519828",
			"supplier_id":    "2da67854-8d6b-4787-a2ce-bde7e07eb1c4",
			"invoice_number": "100-100-100",
			"invoice_date":   "2022-01-01",
		}

		j, err := json.Marshal(form)
		assert.NoError(t, err)

		req, err := http.NewRequest("POST", testUrl, strings.NewReader(string(j)))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()
		s.Router.ServeHTTP(res, req)

		assert.Equal(t, http.StatusCreated, res.Code)

		body, err := io.ReadAll(res.Body)
		assert.NoError(t, err)
		mapBody := make(map[string]any)
		err = json.Unmarshal(body, &mapBody)
		assert.NoError(t, err)

		assert.Equal(t, "Factura creada", mapBody["message"])

		req, err = http.NewRequest("GET", testUrl, nil)
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		res = httptest.NewRecorder()
		s.Router.ServeHTTP(res, req)

		var r []any
		err = json.Unmarshal(res.Body.Bytes(), &r)
		t.Log(r)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(r))
		for _, v := range r {
			assert.Equal(t, "100-100-100", v.(map[string]any)["invoice_number"])
			assert.Equal(t, "2022-01-01T00:00:00Z", v.(map[string]any)["invoice_date"])
			assert.Equal(t, float64(0), v.(map[string]any)["invoice_total"])
			assert.Equal(t, false, v.(map[string]any)["is_balanced"])

			project := v.(map[string]any)["project"].(map[string]any)
			assert.Equal(t, "Project 2", project["name"])
			assert.Equal(t, "1c6020db-39a0-451d-89ee-fdd20d519828", project["id"])

			supplier := v.(map[string]any)["supplier"].(map[string]any)
			assert.Equal(t, "Supplier Name 2", supplier["name"])
			assert.Equal(t, "2da67854-8d6b-4787-a2ce-bde7e07eb1c4", supplier["id"])
		}
	})

}
