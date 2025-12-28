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

func TestApiMaterialsTest(t *testing.T) {
	ctx := context.Background()
	testUrl := "/api/v2/materials"

	pgContainer, err := postgres.Run(ctx,
		"postgres:18-alpine",
		postgres.WithOrderedInitScripts(
			filepath.Join("..", "schema", "tables.sql"),
			filepath.Join("scripts", "seed_categories.sql"),
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

	t.Run("should have no materials", func(t *testing.T) {
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

	t.Run("should be able to create a material", func(t *testing.T) {
		form := map[string]any{
			"name":        "Prueba",
			"code":        "PRU",
			"unit":        "Unidad",
			"category_id": "325a7622-53dd-4ccf-b109-4f98b8daac55",
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

		assert.Equal(t, "Material creado", mapBody["message"])

		req, err = http.NewRequest("GET", testUrl, nil)
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		res = httptest.NewRecorder()
		s.Router.ServeHTTP(res, req)

		var r []any
		err = json.Unmarshal(res.Body.Bytes(), &r)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(r))
		for _, v := range r {
			assert.Equal(t, "Prueba", v.(map[string]any)["name"])
			assert.Equal(t, "PRU", v.(map[string]any)["code"])
			assert.Equal(t, "Unidad", v.(map[string]any)["unit"])
			cat := v.(map[string]any)["category"].(map[string]any)
			assert.Equal(t, "325a7622-53dd-4ccf-b109-4f98b8daac55", cat["id"])
			assert.Equal(t, "prueba", cat["name"])
		}
	})

	t.Run("should show conflict error", func(t *testing.T) {
		form := map[string]any{
			"name":        "Prueba",
			"code":        "PRU",
			"unit":        "Unidad",
			"category_id": "325a7622-53dd-4ccf-b109-4f98b8daac55",
		}

		j, err := json.Marshal(form)
		assert.NoError(t, err)

		req, err := http.NewRequest("POST", testUrl, strings.NewReader(string(j)))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()
		s.Router.ServeHTTP(res, req)

		assert.Equal(t, http.StatusConflict, res.Code)
	})
}
