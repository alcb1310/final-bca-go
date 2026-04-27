package router_test

import (
	"encoding/json"
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

func TestCreateItem(t *testing.T) {
	db := mocks.NewService(t)
	s := router.Router{DB: db}
	assert.NotNil(t, s)
	s.Router()

	server := httptest.NewServer(s.Router())
	defer server.Close()

	testURL := "/api/v2/items"
	testData := []struct {
		name       string
		form       map[string]any
		status     int
		body       map[string]any
		createItem *mocks.Service_CreateItem_Call
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
				"name": "El nombre es obligatorio",
				"code": "El código es obligatorio",
				"unit": "La unidad es obligatoria",
			},
		},
		{
			name: "should pass a code",
			form: map[string]any{
				"name": "Prueba",
			},
			status: http.StatusBadRequest,
			body: map[string]any{
				"code": "El código es obligatorio",
				"unit": "La unidad es obligatoria",
			},
		},
		{
			name: "should pass a unit",
			form: map[string]any{
				"name": "Prueba",
				"code": "prb",
			},
			status: http.StatusBadRequest,
			body: map[string]any{
				"unit": "La unidad es obligatoria",
			},
		},
		{
			name: "should create an item",
			form: map[string]any{
				"name": "Prueba",
				"code": "prb",
				"unit": "m2",
			},
			status: http.StatusCreated,
			body: map[string]any{
				"name": "Prueba",
				"code": "prb",
				"unit": "m2",
				"id":   uuid.Nil.String(),
			},
			createItem: db.EXPECT().CreateItem(&types.Items{
				Name: "Prueba",
				Code: "prb",
				Unit: "m2",
			}).Return(nil),
		},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			var read io.Reader = nil

			if tt.createItem != nil {
				tt.createItem.Times(1)
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
