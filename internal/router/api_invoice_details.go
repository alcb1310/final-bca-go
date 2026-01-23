package router

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/alcb1310/final-bca-go/internal/types"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (rf *Router) GetInvoiceDetails(w http.ResponseWriter, r *http.Request) {
	pId := chi.URLParam(r, "projectId")
	parsedId, err := uuid.Parse(pId)
	inv := types.InvoiceResponse{}
	if err != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		_ = json.NewEncoder(w).Encode(map[string]any{"message": "Id inválido"})
		return
	}

	if inv, err = rf.DB.GetInvoice(parsedId); err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]any{"message": "Factura no encontrada"})
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(err)
		return
	}

	w.WriteHeader(http.StatusNotImplemented)
	_ = json.NewEncoder(w).Encode(inv)
}
