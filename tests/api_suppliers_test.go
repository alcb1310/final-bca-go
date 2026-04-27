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

func TestApiSuppliers(t *testing.T) {
	ctx := context.Background()
	testUrl := "/api/v2/suppliers"
	path := filepath.Join("..", "schema", "tables.sql")

	pgContainer, err := postgres.Run(ctx,
		"postgres:18-alpine",
		postgres.WithInitScripts(path),
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

	s, server, err := createServer(t, ctx, pgContainer)
	assert.NoError(t, err)
	if err != nil {
		return
	}

	defer server.Close()

	t.Run("should have no suppliers", func(t *testing.T) {
		req, err := http.NewRequest("GET", testUrl, nil)
		assert.NoError(t, err)
		res := httptest.NewRecorder()
		s.Router().ServeHTTP(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
		var r []any
		err = json.Unmarshal(res.Body.Bytes(), &r)
		assert.NoError(t, err)
		assert.Equal(t, 0, len(r))
		assert.Equal(t, "[]", strings.TrimSpace(res.Body.String()))
	})

	t.Run("should be able to create a new supplier", func(t *testing.T) {
		form := map[string]any{
			"name":        "Proveedor 1",
			"supplier_id": "1234567890",
		}

		j, err := json.Marshal(form)
		assert.NoError(t, err)

		req, err := http.NewRequest(http.MethodPost, testUrl, strings.NewReader(string(j)))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()
		s.Router().ServeHTTP(res, req)

		assert.Equal(t, http.StatusCreated, res.Code)

		body, err := io.ReadAll(res.Body)
		assert.NoError(t, err)
		mapBody := make(map[string]any)
		err = json.Unmarshal(body, &mapBody)
		assert.NoError(t, err)

		assert.Equal(t, "Proveedor creado correctamente", mapBody["message"])

		req, err = http.NewRequest("GET", testUrl, nil)
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		res = httptest.NewRecorder()
		s.Router().ServeHTTP(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
		var r []any
		err = json.Unmarshal(res.Body.Bytes(), &r)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(r))

		for _, v := range r {
			assert.Equal(t, "Proveedor 1", v.(map[string]any)["name"])
			assert.Equal(t, "1234567890", v.(map[string]any)["supplier_id"])

			name := v.(map[string]any)["contact_name"]
			assert.Equal(t, false, name.(map[string]any)["Valid"])

			email := v.(map[string]any)["contact_email"]
			assert.Equal(t, false, email.(map[string]any)["Valid"])

			phone := v.(map[string]any)["contact_phone"]
			assert.Equal(t, false, phone.(map[string]any)["Valid"])
		}
	})

	t.Run("should be able to create a new supplier", func(t *testing.T) {
		form := map[string]any{
			"name":        "Proveedor 1",
			"supplier_id": "1234567890",
		}

		j, err := json.Marshal(form)
		assert.NoError(t, err)

		req, err := http.NewRequest(http.MethodPost, testUrl, strings.NewReader(string(j)))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()
		s.Router().ServeHTTP(res, req)

		assert.Equal(t, http.StatusConflict, res.Code)
	})

	t.Run("should return not exist when non existent supplier", func(t *testing.T) {
		form := map[string]any{
			"name":        "Proveedor 1",
			"supplier_id": "1234567890",
		}

		j, err := json.Marshal(form)
		assert.NoError(t, err)
		id := uuid.New()

		req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/%s", testUrl, id), strings.NewReader(string(j)))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()
		s.Router().ServeHTTP(res, req)

		assert.Equal(t, http.StatusNotFound, res.Code)
		body, err := io.ReadAll(res.Body)
		assert.NoError(t, err)
		mapBody := make(map[string]any)
		err = json.Unmarshal(body, &mapBody)
		assert.NoError(t, err)

		assert.Equal(t, "Proveedor no encontrado", mapBody["message"])
	})

	t.Run("should update a supplier", func(t *testing.T) {
		req, err := http.NewRequest("GET", testUrl, nil)
		assert.NoError(t, err)
		res := httptest.NewRecorder()
		s.Router().ServeHTTP(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
		var r []any
		err = json.Unmarshal(res.Body.Bytes(), &r)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(r))

		id := r[0].(map[string]any)["id"].(string)

		form := map[string]any{
			"name":          "Proveedor 1 Edited",
			"supplier_id":   "123456789001",
			"contact_name":  "Andres",
			"contact_phone": "1234567",
			"contact_email": "Yk6bM@example.com",
		}

		j, err := json.Marshal(form)
		assert.NoError(t, err)

		req, err = http.NewRequest(http.MethodPut, fmt.Sprintf("%s/%s", testUrl, id), strings.NewReader(string(j)))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		res = httptest.NewRecorder()
		s.Router().ServeHTTP(res, req)

		assert.Equal(t, http.StatusNoContent, res.Code)

		req, err = http.NewRequest("GET", testUrl, nil)
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		res = httptest.NewRecorder()
		s.Router().ServeHTTP(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
		err = json.Unmarshal(res.Body.Bytes(), &r)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(r))

		for _, v := range r {
			assert.Equal(t, "Proveedor 1 Edited", v.(map[string]any)["name"])
			assert.Equal(t, "123456789001", v.(map[string]any)["supplier_id"])

			name := v.(map[string]any)["contact_name"]
			assert.Equal(t, true, name.(map[string]any)["Valid"])
			assert.Equal(t, "Andres", name.(map[string]any)["String"])

			email := v.(map[string]any)["contact_email"]
			assert.Equal(t, true, email.(map[string]any)["Valid"])
			assert.Equal(t, "Yk6bM@example.com", email.(map[string]any)["String"])

			phone := v.(map[string]any)["contact_phone"]
			assert.Equal(t, true, phone.(map[string]any)["Valid"])
			assert.Equal(t, "1234567", phone.(map[string]any)["String"])
		}
	})
}
