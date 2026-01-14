package router

import (
	"encoding/json"
	"net/http"

	"github.com/alcb1310/final-bca-go/internal/types"
)

func (rf *Router) GetInvoices(w http.ResponseWriter, r *http.Request) {
	invoices := []types.InvoiceResponse{}
	var err error

	if invoices, err = rf.DB.GetInvoices(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(invoices)
}

func (rf *Router) CreateInvoice(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}
