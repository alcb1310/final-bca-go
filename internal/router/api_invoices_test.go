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

func TestApiInvoices(t *testing.T) {
	db := mocks.NewService(t)
	s := router.NewRouter(db)
	s.GenerateRoutes()
	testURL := "/api/v2/invoices"
	testData := []struct {
		name     string
		form     map[string]any
		status   int
		body     map[string]any
		project  *mocks.Service_GetProject_Call
		supplier *mocks.Service_GetSupplier_Call
		create   *mocks.Service_CreateInvoice_Call
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
			name:   "should pass a body",
			form:   map[string]any{},
			status: http.StatusBadRequest,
			body: map[string]any{
				"project_id":     "El proyecto es obligatorio",
				"supplier_id":    "El proveedor es obligatorio",
				"invoice_number": "El numero de factura es obligatorio",
				"invoice_date":   "La fecha de la factura es obligatoria",
			},
		},
		{
			name: "should pass a non empty project id",
			form: map[string]any{
				"project_id": "",
			},
			status: http.StatusBadRequest,
			body: map[string]any{
				"project_id":     "El código del proyecto es inválido",
				"supplier_id":    "El proveedor es obligatorio",
				"invoice_number": "El numero de factura es obligatorio",
				"invoice_date":   "La fecha de la factura es obligatoria",
			},
		},
		{
			name: "should pass a valid project id",
			form: map[string]any{
				"project_id": "invalid",
			},
			status: http.StatusBadRequest,
			body: map[string]any{
				"project_id":     "El código del proyecto es inválido",
				"supplier_id":    "El proveedor es obligatorio",
				"invoice_number": "El numero de factura es obligatorio",
				"invoice_date":   "La fecha de la factura es obligatoria",
			},
		},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			var read io.Reader = nil

			if tt.project != nil {
				tt.project.Times(1)
			}
			if tt.create != nil {
				tt.create.Times(1)
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
