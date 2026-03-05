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
	"github.com/alcb1310/final-bca-go/mocks"
	"github.com/stretchr/testify/assert"
)

func TestApiInvoiceDetails(t *testing.T) {
	// invoiceId := uuid.New()
	db := mocks.NewService(t)
	s := router.NewRouter(db)
	s.GenerateRoutes()

	testData := []struct {
		name      string
		form      map[string]any
		body      map[string]any
		status    int
		invoiceId string
		invoice   *mocks.Service_GetInvoice_Call
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
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			testURL := fmt.Sprintf("/api/v2/invoices/%s/details", tt.invoiceId)
			var read io.Reader = nil
			// invoice := db.EXPECT().GetInvoice(tt.invoiceId).Return(types.InvoiceResponse{Id: tt.invoiceId}, nil)
			if tt.invoice != nil {
				tt.invoice.Times(1)
			}

			if tt.form != nil {
				data, err := json.Marshal(tt.form)
				assert.NoError(t, err)
				read = strings.NewReader(string(data))
			}

			req, err := http.NewRequest("GET", testURL, read)
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			res := httptest.NewRecorder()
			s.Router.ServeHTTP(res, req)
			assert.Equal(t, tt.status, res.Code)

			body, err := io.ReadAll(res.Body)
			t.Log(string(body))
			assert.NoError(t, err)

			mapBody := make(map[string]any)
			err = json.Unmarshal(body, &mapBody)
			t.Log(mapBody)
			assert.NoError(t, err)

			assert.Equal(t, len(tt.body), len(mapBody))

			for k, v := range tt.body {
				assert.Equal(t, v, mapBody[k])
			}
		})
	}
}
