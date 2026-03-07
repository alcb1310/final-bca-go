package router

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/alcb1310/final-bca-go/internal/types"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
)

func (rf *Router) GetInvoiceDetails(w http.ResponseWriter, r *http.Request) {
	pId := chi.URLParam(r, "projectId")
	parsedId, err := uuid.Parse(pId)
	if err != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		_ = json.NewEncoder(w).Encode(map[string]any{"message": "Id inválido"})
		return
	}

	if _, err = rf.DB.GetInvoice(parsedId); err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]any{"message": "Factura no encontrada"})
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(err)
		return
	}

	invoiceDetails, err := rf.DB.GetInvoiceDetails(parsedId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(invoiceDetails)
}

func (rf *Router) CreateInvoiceDetail(w http.ResponseWriter, r *http.Request) {
	pId := chi.URLParam(r, "projectId")
	parsedId, err := uuid.Parse(pId)
	if err != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		_ = json.NewEncoder(w).Encode(map[string]any{"message": "Id inválido"})
		return
	}

	if _, err = rf.DB.GetInvoice(parsedId); err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]any{"message": "Factura no encontrada"})
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(err)
		return
	}

	p := make(map[string]any)
	invoiceDetail := types.InvoiceDetailsCreate{
		InvoiceId: parsedId,
	}
	errorResponse := make(map[string]any)
	var ok bool

	if r.Body == http.NoBody || r.Body == nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		errorResponse["message"] = "Falta el cuerpo de la solicitud"
		_ = json.NewEncoder(w).Encode(errorResponse)
		return
	}

	if err = json.NewDecoder(r.Body).Decode(&p); err != nil {
		errorResponse["message"] = "Cuerpo de la solicitud no válido"
	}

	if val, ok := p["budget_item_id"].(string); !ok {
		errorResponse["budget_item_id"] = "El id de la partida es obligatorio"
	} else {
		invoiceDetail.BudgetItemId, err = uuid.Parse(val)
		if err != nil {
			errorResponse["budget_item_id"] = "El id de la partida es inválido"
		}
	}

	if invoiceDetail.Quantity, ok = p["quantity"].(float64); !ok {
		errorResponse["quantity"] = "La cantidad es obligatoria"
	}

	if invoiceDetail.Cost, ok = p["cost"].(float64); !ok {
		errorResponse["cost"] = "El costo es obligatorio"
	}

	if len(errorResponse) > 0 {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errorResponse)
		return
	}

	if err = rf.DB.CreateInvoiceDetail(invoiceDetail); err != nil {
		var e *pgconn.PgError
		if errors.As(err, &e) {
			switch e.Code {
			case "23505":
				w.WriteHeader(http.StatusConflict)
				_ = json.NewEncoder(w).Encode(map[string]any{"message": "El detalle de la factura ya existe"})
				return
			}
		}
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]any{"message": "Detalle de factura creado"})
}

func (rf *Router) DeleteInvoiceDetail(w http.ResponseWriter, r *http.Request) {
	pId := chi.URLParam(r, "projectId")
	parsedId, err := uuid.Parse(pId)
	if err != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		_ = json.NewEncoder(w).Encode(map[string]any{"message": "Id inválido"})
		return
	}

	if _, err = rf.DB.GetInvoice(parsedId); err != nil {
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
}
