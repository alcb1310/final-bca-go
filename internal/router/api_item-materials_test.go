package router_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/alcb1310/final-bca-go/internal/router"
	"github.com/alcb1310/final-bca-go/internal/types"
	"github.com/alcb1310/final-bca-go/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateItemMaterial(t *testing.T) {
	db := mocks.NewService(t)
	s := router.NewRouter(db)
	s.GenerateRoutes()
	rubroId := uuid.New()
	testURL := fmt.Sprintf("/api/v2/items/%s/materials", rubroId)
	itemExpect := db.EXPECT().GetItem(rubroId).Return(types.Items{
		Id:   rubroId,
		Name: "test",
		Code: "test",
		Unit: "test",
	}, nil)

	testData := []struct {
		name    string
		form    map[string]any
		status  int
		body    map[string]any
		getItem *mocks.Service_GetItem_Call
	}{
		{
			name:   "should pass a form",
			form:   nil,
			status: http.StatusUnprocessableEntity,
			body: map[string]any{
				"message": "Falta el cuerpo de la solicitud",
			},
			getItem: itemExpect,
		},
		{
			name:   "should pass a material_id",
			form:   map[string]any{},
			status: http.StatusBadRequest,
			body: map[string]any{
				"material_id": "El material es obligatorio",
				"quantity":    "La cantidad es obligatoria",
			},
			getItem: itemExpect,
		},
		{
			name: "should pass a non empty material_id",
			form: map[string]any{
				"material_id": "",
			},
			status: http.StatusBadRequest,
			body: map[string]any{
				"material_id": "El material es invalido",
				"quantity":    "La cantidad es obligatoria",
			},
			getItem: itemExpect,
		},
		{
			name: "should pass a valid material_id",
			form: map[string]any{
				"material_id": "invalid",
			},
			status: http.StatusBadRequest,
			body: map[string]any{
				"material_id": "El material es invalido",
				"quantity":    "La cantidad es obligatoria",
			},
			getItem: itemExpect,
		},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			var read io.Reader = nil

			if tt.getItem != nil {
				tt.getItem.Times(1)
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
