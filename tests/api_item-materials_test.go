package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
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

	t.Run("should have no rubro materials", func(t *testing.T) {
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

	t.Run("should not find rubro on POST", func(t *testing.T) {
		invalidUUID := uuid.New()
		testUrl := fmt.Sprintf("/api/v2/items/%s/materials", invalidUUID.String())
		form := map[string]any{
			"material_id": "b3fba400-acad-40a6-9ca3-17871151bc0f",
			"quantity":    32.5,
		}

		j, err := json.Marshal(form)
		assert.NoError(t, err)

		req, err := http.NewRequest("POST", testUrl, strings.NewReader(string(j)))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()
		s.Router.ServeHTTP(res, req)

		assert.Equal(t, http.StatusNotFound, res.Code)

		body, err := io.ReadAll(res.Body)
		assert.NoError(t, err)
		mapBody := make(map[string]any)
		err = json.Unmarshal(body, &mapBody)
		assert.NoError(t, err)

		assert.Equal(t, "Rubro no encontrado", mapBody["message"])
	})

	t.Run("should be able to create a rubro material", func(t *testing.T) {
		form := map[string]any{
			"material_id": "b3fba400-acad-40a6-9ca3-17871151bc0f",
			"quantity":    32.5,
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

		t.Log(mapBody)
		assert.Equal(t, "b3fba400-acad-40a6-9ca3-17871151bc0f", mapBody["material_id"])
		assert.Equal(t, "2d257121-43e8-4b00-947d-b05fa54b36ac", mapBody["item_id"])
		assert.Equal(t, 32.5, mapBody["quantity"])

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
			assert.Equal(t, "b3fba400-acad-40a6-9ca3-17871151bc0f", v.(map[string]any)["material_id"])
			assert.Equal(t, "prueba material", v.(map[string]any)["material_name"])
			assert.Equal(t, "pm123", v.(map[string]any)["material_code"])
			assert.Equal(t, "pri", v.(map[string]any)["material_unit"])
			assert.Equal(t, "2d257121-43e8-4b00-947d-b05fa54b36ac", v.(map[string]any)["item_id"])
			assert.Equal(t, "prueba item", v.(map[string]any)["item_name"])
			assert.Equal(t, "pr123", v.(map[string]any)["item_code"])
			assert.Equal(t, "pri", v.(map[string]any)["item_unit"])
			assert.Equal(t, 32.5, v.(map[string]any)["quantity"])
		}
	})

	t.Run("should show conflict error", func(t *testing.T) {
		form := map[string]any{
			"material_id": "b3fba400-acad-40a6-9ca3-17871151bc0f",
			"quantity":    32.5,
		}

		j, err := json.Marshal(form)
		assert.NoError(t, err)

		req, err := http.NewRequest("POST", testUrl, strings.NewReader(string(j)))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()
		s.Router.ServeHTTP(res, req)

		assert.Equal(t, http.StatusConflict, res.Code)

		body, err := io.ReadAll(res.Body)
		assert.NoError(t, err)
		mapBody := make(map[string]any)
		err = json.Unmarshal(body, &mapBody)
		assert.NoError(t, err)

		assert.Equal(t, "Material ya asociado al rubro", mapBody["message"])
	})

	t.Run("should not find rubro on PUT", func(t *testing.T) {
		invalidUUID := uuid.New()
		materialID := "b3fba400-acad-40a6-9ca3-17871151bc0f"
		testUrl := fmt.Sprintf("/api/v2/items/%s/materials/%s", invalidUUID.String(), materialID)
		form := map[string]any{
			"material_id": "b3fba400-acad-40a6-9ca3-17871151bc0f",
			"quantity":    32.5,
		}

		j, err := json.Marshal(form)
		assert.NoError(t, err)

		req, err := http.NewRequest("PUT", testUrl, strings.NewReader(string(j)))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()
		s.Router.ServeHTTP(res, req)

		assert.Equal(t, http.StatusNotFound, res.Code)

		body, err := io.ReadAll(res.Body)
		assert.NoError(t, err)
		mapBody := make(map[string]any)
		err = json.Unmarshal(body, &mapBody)
		assert.NoError(t, err)

		assert.Equal(t, "Rubro no encontrado", mapBody["message"])
	})
}
