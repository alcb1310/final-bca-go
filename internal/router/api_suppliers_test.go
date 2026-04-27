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
	"github.com/stretchr/testify/assert"
)

func TestApiCreateSupplier(t *testing.T) {
	db := mocks.NewService(t)
	s := router.Router{DB: db}
	assert.NotNil(t, s)
	s.Router()

	server := httptest.NewServer(s.Router())
	defer server.Close()

	testURL := "/api/v2/suppliers"
	testData := []struct {
		name           string
		form           map[string]any
		status         int
		body           map[string]any
		createSupplier *mocks.Service_CreateSupplier_Call
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
				"supplier_id": "El RUC es obligatorio",
			},
		},
		{
			name: "should pass a name",
			form: map[string]any{
				"name": "",
			},
			status: http.StatusBadRequest,
			body: map[string]any{
				"name":        "El nombre es obligatorio",
				"supplier_id": "El RUC es obligatorio",
			},
		},
		{
			name: "should pass a name",
			form: map[string]any{
				"name": "Test Supplier",
			},
			status: http.StatusBadRequest,
			body: map[string]any{
				"supplier_id": "El RUC es obligatorio",
			},
		},
		{
			name: "should pass a name",
			form: map[string]any{
				"name":        "Test Supplier",
				"supplier_id": "",
			},
			status: http.StatusBadRequest,
			body: map[string]any{
				"supplier_id": "El RUC es obligatorio",
			},
		},
		{
			name: "should pass a name",
			form: map[string]any{
				"name":        "Test Supplier",
				"supplier_id": "1234567890",
			},
			status: http.StatusCreated,
			body: map[string]any{
				"message": "Proveedor creado correctamente",
			},
			createSupplier: db.EXPECT().CreateSupplier(types.Supplier{
				Name:       "Test Supplier",
				SupplierId: "1234567890",
			}).Return(nil),
		},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			var read io.Reader = nil

			if tt.createSupplier != nil {
				tt.createSupplier.Times(1)
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

			for k, v := range tt.body {
				assert.Equal(t, v, mapBody[k])
			}
		})
	}
}
