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
	s := router.Router{DB: db}
	assert.NotNil(t, s)
	s.Router()

	server := httptest.NewServer(s.Router())
	defer server.Close()

	rubroId := uuid.New()
	testURL := fmt.Sprintf("/api/v2/items/%s/materials", rubroId)
	itemExpect := db.EXPECT().GetItem(rubroId).Return(types.Items{
		Id:   rubroId,
		Name: "test",
		Code: "test",
		Unit: "test",
	}, nil)
	materialId := uuid.New()

	testData := []struct {
		name               string
		form               map[string]any
		status             int
		body               map[string]any
		getItem            *mocks.Service_GetItem_Call
		createItemMaterial *mocks.Service_CreateItemMaterial_Call
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
		{
			name: "should pass a quantity",
			form: map[string]any{
				"material_id": materialId.String(),
			},
			status: http.StatusBadRequest,
			body: map[string]any{
				"quantity": "La cantidad es obligatoria",
			},
			getItem: itemExpect,
		},
		{
			name: "should pass a numeric quantity",
			form: map[string]any{
				"material_id": materialId.String(),
				"quantity":    "invalid",
			},
			status: http.StatusBadRequest,
			body: map[string]any{
				"quantity": "La cantidad es obligatoria",
			},
			getItem: itemExpect,
		},
		{
			name: "should create a new item material",
			form: map[string]any{
				"material_id": materialId.String(),
				"quantity":    45.32,
			},
			status: http.StatusCreated,
			body: map[string]any{
				"item_id":     rubroId.String(),
				"material_id": materialId.String(),
				"quantity":    45.32,
			},
			getItem: itemExpect,
			createItemMaterial: db.EXPECT().CreateItemMaterial(types.ItemMaterialCreate{
				ItemId:     rubroId,
				MaterialId: materialId,
				Quantity:   45.32,
			}).Return(nil),
		},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			var read io.Reader = nil

			if tt.getItem != nil {
				tt.getItem.Times(1)
			}

			if tt.createItemMaterial != nil {
				tt.createItemMaterial.Times(1)
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
			s.Router().ServeHTTP(res, req)
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
