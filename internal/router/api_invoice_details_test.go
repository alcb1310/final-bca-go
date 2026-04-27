package router_test

import (
	"database/sql"
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

func TestApiInvoiceDetails(t *testing.T) {
	db := mocks.NewService(t)
	s := router.Router{DB: db}
	assert.NotNil(t, s)
	s.Router()

	server := httptest.NewServer(s.Router())
	defer server.Close()

	invoiceId := uuid.New()
	budgetItemId := uuid.New()

	testData := []struct {
		name          string
		form          map[string]any
		body          map[string]any
		status        int
		invoiceId     string
		invoice       *mocks.Service_GetInvoice_Call
		createInvoice *mocks.Service_CreateInvoiceDetail_Call
	}{
		{
			name:      "should pass a valid invoice id",
			invoiceId: "invalid",
			invoice:   nil,
			form:      make(map[string]any),
			status:    http.StatusNotAcceptable,
			body: map[string]any{
				"message": "Id inválido",
			},
		},
		{
			name:      "should return not found",
			invoiceId: invoiceId.String(),
			invoice:   db.EXPECT().GetInvoice(invoiceId).Return(types.InvoiceResponse{}, sql.ErrNoRows),
			form:      make(map[string]any),
			status:    http.StatusNotFound,
			body: map[string]any{
				"message": "Factura no encontrada",
			},
		},
		{
			name:      "should pass a budget item",
			invoiceId: invoiceId.String(),
			invoice:   db.EXPECT().GetInvoice(invoiceId).Return(types.InvoiceResponse{}, nil),
			form:      make(map[string]any),
			status:    http.StatusBadRequest,
			body: map[string]any{
				"quantity":       "La cantidad es obligatoria",
				"cost":           "El costo es obligatorio",
				"budget_item_id": "El id de la partida es obligatorio",
			},
		},
		{
			name:      "should pass a valid budget item",
			invoiceId: invoiceId.String(),
			invoice:   db.EXPECT().GetInvoice(invoiceId).Return(types.InvoiceResponse{}, nil),
			form: map[string]any{
				"budget_item_id": "invalid",
			},
			status: http.StatusBadRequest,
			body: map[string]any{
				"quantity":       "La cantidad es obligatoria",
				"cost":           "El costo es obligatorio",
				"budget_item_id": "El id de la partida es inválido",
			},
		},
		{
			name:      "should pass a quantity",
			invoiceId: invoiceId.String(),
			invoice:   db.EXPECT().GetInvoice(invoiceId).Return(types.InvoiceResponse{}, nil),
			form: map[string]any{
				"budget_item_id": budgetItemId.String(),
			},
			status: http.StatusBadRequest,
			body: map[string]any{
				"quantity": "La cantidad es obligatoria",
				"cost":     "El costo es obligatorio",
			},
		},
		{
			name:      "should pass a valid quantity",
			invoiceId: invoiceId.String(),
			invoice:   db.EXPECT().GetInvoice(invoiceId).Return(types.InvoiceResponse{}, nil),
			form: map[string]any{
				"budget_item_id": budgetItemId.String(),
				"quantity":       "invalid",
			},
			status: http.StatusBadRequest,
			body: map[string]any{
				"quantity": "La cantidad es obligatoria",
				"cost":     "El costo es obligatorio",
			},
		},
		{
			name:      "should pass a cost",
			invoiceId: invoiceId.String(),
			invoice:   db.EXPECT().GetInvoice(invoiceId).Return(types.InvoiceResponse{}, nil),
			form: map[string]any{
				"budget_item_id": budgetItemId.String(),
				"quantity":       10.5,
			},
			status: http.StatusBadRequest,
			body: map[string]any{
				"cost": "El costo es obligatorio",
			},
		},
		{
			name:      "should pass a valid cost",
			invoiceId: invoiceId.String(),
			invoice:   db.EXPECT().GetInvoice(invoiceId).Return(types.InvoiceResponse{}, nil),
			form: map[string]any{
				"budget_item_id": budgetItemId.String(),
				"quantity":       10.5,
				"cost":           "invalid",
			},
			status: http.StatusBadRequest,
			body: map[string]any{
				"cost": "El costo es obligatorio",
			},
		},
		{
			name:      "should create the invoice detail",
			invoiceId: invoiceId.String(),
			invoice:   db.EXPECT().GetInvoice(invoiceId).Return(types.InvoiceResponse{}, nil),
			form: map[string]any{
				"budget_item_id": budgetItemId.String(),
				"quantity":       10.5,
				"cost":           5,
			},
			status: http.StatusCreated,
			body: map[string]any{
				"message": "Detalle de factura creado",
			},
			createInvoice: db.EXPECT().CreateInvoiceDetail(types.InvoiceDetailsCreate{
				InvoiceId:    invoiceId,
				BudgetItemId: budgetItemId,
				Quantity:     10.5,
				Cost:         5,
			}).Return(nil),
		},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			testURL := fmt.Sprintf("/api/v2/invoices/%s/details", tt.invoiceId)
			var read io.Reader = nil
			if tt.invoice != nil {
				tt.invoice.Times(1)
			}

			if tt.form != nil {
				data, err := json.Marshal(tt.form)
				assert.NoError(t, err)
				read = strings.NewReader(string(data))
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
