package router_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/alcb1310/final-bca-go/internal/router"
	"github.com/alcb1310/final-bca-go/mocks"
	"github.com/stretchr/testify/assert"
)

func TestCreateMaterial(t *testing.T) {
	db := mocks.NewService(t)
	s := router.NewRouter(db)
	s.GenerateRoutes()
	testURL := "/api/v2/materials"
	testData := []struct {
		name           string
		form           map[string]any
		status         int
		body           map[string]any
		createMaterial *mocks.Service_CreateMaterial_Call
	}{
		{
			name:   "should pass a form",
			form:   nil,
			status: http.StatusUnprocessableEntity,
			body: map[string]any{
				"message": "Falta el cuerpo de la solicitud",
			},
		},
		{
			name:   "should pass a name",
			form:   map[string]any{},
			status: http.StatusBadRequest,
			body: map[string]any{
				"name":        "El nombre es obligatorio",
				"code":        "El código es obligatorio",
				"unit":        "La unidad es obligatoria",
				"category_id": "La categoría es obligatoria",
			},
		},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			var read io.Reader = nil

			if tt.createMaterial != nil {
				tt.createMaterial.Times(1)
			}

			if tt.form != nil {
				j, err := json.Marshal(tt.form)
				assert.NoError(t, err)
				read = strings.NewReader(string(j))
			}

			req, err := http.NewRequest(http.MethodPost, testURL, read)
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			res := httptest.NewRecorder()
			s.Router.ServeHTTP(res, req)
			assert.Equal(t, tt.status, res.Code)

			body, err := io.ReadAll(res.Body)
			assert.NoError(t, err)

			mapBody := make(map[string]any)
			err = json.Unmarshal(body, &mapBody)
			assert.NoError(t, err)

			assert.Equal(t, len(tt.body), len(mapBody))

			for k, v := range tt.body {
				assert.Equal(t, v, mapBody[k])
			}
		})
	}
}
