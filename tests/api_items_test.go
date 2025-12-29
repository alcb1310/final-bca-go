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

func TestApiItemsTest(t *testing.T) {
	ctx := context.Background()
	testUrl := "/api/v2/items"
	t.Log(testUrl)
	path := filepath.Join("..", "schema", "tables.sql")

	pgContainer, err := postgres.Run(ctx,
		"postgres:18-alpine",
		postgres.WithOrderedInitScripts(path),
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

	t.Run("should have no rubros", func(t *testing.T) {
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

	t.Run("should be able to create a rubro", func(t *testing.T) {
		form := map[string]any{
			"name": "Prueba",
			"code": "PRU",
			"unit": "Unidad",
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
		assert.Equal(t, "Prueba", mapBody["name"])
		assert.Equal(t, "PRU", mapBody["code"])
		assert.Equal(t, "Unidad", mapBody["unit"])
		id := mapBody["id"].(string)
		assert.NotEmpty(t, id)

		_, err = uuid.Parse(id)
		assert.NoError(t, err)

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
		}
	})

	t.Run("should show conflict error", func(t *testing.T) {
		form := map[string]any{
			"name": "Prueba",
			"code": "PRU",
			"unit": "Unidad",
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

		assert.Equal(t, "El rubro ya existe", mapBody["message"])
	})

	t.Run("should return not exist when non existent rubro", func(t *testing.T) {
		form := map[string]any{
			"code": "100",
			"name": "Test",
			"unit": "u",
		}

		j, err := json.Marshal(form)
		assert.NoError(t, err)
		id := uuid.New()

		req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/%s", testUrl, id), strings.NewReader(string(j)))
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

	t.Run("should be able to update a rubro", func(t *testing.T) {
		req, err := http.NewRequest("GET", testUrl, nil)
		assert.NoError(t, err)
		res := httptest.NewRecorder()
		s.Router.ServeHTTP(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
		var r []any
		err = json.Unmarshal(res.Body.Bytes(), &r)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(r))

		id := r[0].(map[string]any)["id"].(string)

		form := map[string]any{
			"code": "100",
			"name": "Test",
			"unit": "u",
		}

		j, err := json.Marshal(form)
		assert.NoError(t, err)

		req, err = http.NewRequest(http.MethodPut, fmt.Sprintf("%s/%s", testUrl, id), strings.NewReader(string(j)))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		res = httptest.NewRecorder()
		s.Router.ServeHTTP(res, req)

		assert.Equal(t, http.StatusNoContent, res.Code)

		req, err = http.NewRequest("GET", testUrl, nil)
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		res = httptest.NewRecorder()
		s.Router.ServeHTTP(res, req)

		err = json.Unmarshal(res.Body.Bytes(), &r)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(r))
		for _, v := range r {
			assert.Equal(t, "Test", v.(map[string]any)["name"])
			assert.Equal(t, "100", v.(map[string]any)["code"])
			assert.Equal(t, "u", v.(map[string]any)["unit"])
		}
	})
}
