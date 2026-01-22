package router_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/alcb1310/final-bca-go/internal/router"
	"github.com/alcb1310/final-bca-go/internal/types"
	"github.com/alcb1310/final-bca-go/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestApiInvoices(t *testing.T) {
	projectId := uuid.New()
	supplierId := uuid.New()
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
		{
			name: "should pass a supplier id",
			form: map[string]any{
				"project_id": projectId.String(),
			},
			status: http.StatusBadRequest,
			body: map[string]any{
				"supplier_id":    "El proveedor es obligatorio",
				"invoice_number": "El numero de factura es obligatorio",
				"invoice_date":   "La fecha de la factura es obligatoria",
			},
			project: db.EXPECT().GetProject(projectId).Return(types.Project{}, nil),
		},
		{
			name: "should pass a non empty supplier id",
			form: map[string]any{
				"project_id":  projectId.String(),
				"supplier_id": "",
			},
			status: http.StatusBadRequest,
			body: map[string]any{
				"supplier_id":    "El código del proveedor es inválido",
				"invoice_number": "El numero de factura es obligatorio",
				"invoice_date":   "La fecha de la factura es obligatoria",
			},
			project: db.EXPECT().GetProject(projectId).Return(types.Project{}, nil),
		},
		{
			name: "should pass a valid supplier id",
			form: map[string]any{
				"project_id":  projectId.String(),
				"supplier_id": "invalid",
			},
			status: http.StatusBadRequest,
			body: map[string]any{
				"supplier_id":    "El código del proveedor es inválido",
				"invoice_number": "El numero de factura es obligatorio",
				"invoice_date":   "La fecha de la factura es obligatoria",
			},
			project: db.EXPECT().GetProject(projectId).Return(types.Project{}, nil),
		},
		{
			name: "should pass an invoice number",
			form: map[string]any{
				"project_id":  projectId.String(),
				"supplier_id": supplierId.String(),
			},
			status: http.StatusBadRequest,
			body: map[string]any{
				"invoice_number": "El numero de factura es obligatorio",
				"invoice_date":   "La fecha de la factura es obligatoria",
			},
			project:  db.EXPECT().GetProject(projectId).Return(types.Project{}, nil),
			supplier: db.EXPECT().GetSupplier(supplierId).Return(types.Supplier{}, nil),
		},
		{
			name: "should pass a non empty invoice number",
			form: map[string]any{
				"project_id":     projectId.String(),
				"supplier_id":    supplierId.String(),
				"invoice_number": "",
			},
			status: http.StatusBadRequest,
			body: map[string]any{
				"invoice_number": "El numero de factura es obligatorio",
				"invoice_date":   "La fecha de la factura es obligatoria",
			},
			project:  db.EXPECT().GetProject(projectId).Return(types.Project{}, nil),
			supplier: db.EXPECT().GetSupplier(supplierId).Return(types.Supplier{}, nil),
		},
		{
			name: "should pass an invoice date",
			form: map[string]any{
				"project_id":     projectId.String(),
				"supplier_id":    supplierId.String(),
				"invoice_number": "001-001-100",
			},
			status: http.StatusBadRequest,
			body: map[string]any{
				"invoice_date": "La fecha de la factura es obligatoria",
			},
			project:  db.EXPECT().GetProject(projectId).Return(types.Project{}, nil),
			supplier: db.EXPECT().GetSupplier(supplierId).Return(types.Supplier{}, nil),
		},
		{
			name: "should pass a non empty invoice date",
			form: map[string]any{
				"project_id":     projectId.String(),
				"supplier_id":    supplierId.String(),
				"invoice_number": "001-001-100",
				"invoice_date":   "",
			},
			status: http.StatusBadRequest,
			body: map[string]any{
				"invoice_date": "La fecha de la factura es inválida",
			},
			project:  db.EXPECT().GetProject(projectId).Return(types.Project{}, nil),
			supplier: db.EXPECT().GetSupplier(supplierId).Return(types.Supplier{}, nil),
		},
		{
			name: "should pass a valid invoice date",
			form: map[string]any{
				"project_id":     projectId.String(),
				"supplier_id":    supplierId.String(),
				"invoice_number": "001-001-100",
				"invoice_date":   "2026-14-14",
			},
			status: http.StatusBadRequest,
			body: map[string]any{
				"invoice_date": "La fecha de la factura es inválida",
			},
			project:  db.EXPECT().GetProject(projectId).Return(types.Project{}, nil),
			supplier: db.EXPECT().GetSupplier(supplierId).Return(types.Supplier{}, nil),
		},
		{
			name: "should create an invoice",
			form: map[string]any{
				"project_id":     projectId.String(),
				"supplier_id":    supplierId.String(),
				"invoice_number": "001-001-100",
				"invoice_date":   "2026-12-14",
			},
			status: http.StatusCreated,
			body: map[string]any{
				"project_id":     projectId.String(),
				"supplier_id":    supplierId.String(),
				"invoice_number": "001-001-100",
				"invoice_date":   "2026-12-14",
				"invoice_total":  float64(0.00),
				"is_balanced":    false,
				"id":             "00000000-0000-0000-0000-000000000000",
			},
			project:  db.EXPECT().GetProject(projectId).Return(types.Project{}, nil),
			supplier: db.EXPECT().GetSupplier(supplierId).Return(types.Supplier{}, nil),
			create: db.EXPECT().CreateInvoice(&types.InvoiceCreate{
				ProjectId:     projectId,
				SupplierId:    supplierId,
				InvoiceNumber: "001-001-100",
				InvoiceDate:   time.Date(2026, 12, 14, 0, 0, 0, 0, time.Now().Location()),
				InvoiceTotal:  0.00,
				IsBalanced:    false,
			}).Return(nil),
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
				if k == "invoice_date" {
					assert.Contains(t, mapBody[k], v.(string))
					continue
				}
				assert.Equal(t, v, mapBody[k])
			}
		})
	}
}
