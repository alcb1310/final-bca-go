package router

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/alcb1310/final-bca-go/internal/types"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
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
	errorResponse := make(map[string]any)
	p := make(map[string]any)
	var err error
	invoice := types.InvoiceCreate{
		IsBalanced:   false,
		InvoiceTotal: 0,
	}
	var val string
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

	if val, ok = p["project_id"].(string); !ok {
		errorResponse["project_id"] = "El proyecto es obligatorio"
	} else {
		invoice.ProjectId, err = uuid.Parse(val)
		if err != nil {
			errorResponse["project_id"] = "El código del proyecto es inválido"
		} else {
			if _, err = rf.DB.GetProject(invoice.ProjectId); err != nil {
				w.WriteHeader(http.StatusNotFound)
				errorResponse["message"] = "El proyecto no existe"
				_ = json.NewEncoder(w).Encode(errorResponse)
				return
			}
		}
	}

	if val, ok = p["supplier_id"].(string); !ok {
		errorResponse["supplier_id"] = "El proveedor es obligatorio"
	} else {
		invoice.SupplierId, err = uuid.Parse(val)
		if err != nil {
			errorResponse["supplier_id"] = "El código del proveedor es inválido"
		} else {
			if _, err = rf.DB.GetSupplier(invoice.SupplierId); err != nil {
				w.WriteHeader(http.StatusNotFound)
				errorResponse["message"] = "El proveedor no existe"
				_ = json.NewEncoder(w).Encode(errorResponse)
				return
			}
		}
	}

	if invoice.InvoiceNumber, ok = p["invoice_number"].(string); !ok {
		errorResponse["invoice_number"] = "El numero de factura es obligatorio"
	} else if len(invoice.InvoiceNumber) == 0 {
		errorResponse["invoice_number"] = "El numero de factura es obligatorio"
	}

	if val, ok = p["invoice_date"].(string); !ok {
		errorResponse["invoice_date"] = "La fecha de la factura es obligatoria"
	} else {
		invoice.InvoiceDate, err = time.ParseInLocation("2006-01-02", val, time.Local)
		if err != nil {
			errorResponse["invoice_date"] = "La fecha de la factura es inválida"
		}
	}

	if len(errorResponse) > 0 {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errorResponse)
		return
	}

	if err = rf.DB.CreateInvoice(invoice); err != nil {
		var e *pgconn.PgError
		if errors.As(err, &e) {
			if e.Code == "23505" {
				w.WriteHeader(http.StatusConflict)
				_ = json.NewEncoder(w).Encode(err)
				return
			}
		}
		slog.Error("CreateInvoice:", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]any{"message": "Factura creada"})
}

func (rr *Router) DeleteInvoice(w http.ResponseWriter, r *http.Request) {
	pId := chi.URLParam(r, "projectId")
	parsedId, err := uuid.Parse(pId)
	inv := types.InvoiceResponse{}
	if err != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		_ = json.NewEncoder(w).Encode(map[string]any{"message": "Id inválido"})
		return
	}

	if inv, err = rr.DB.GetInvoice(parsedId); err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]any{"message": "Proyecto no encontrado"})
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(err)
		return
	}

	if inv.InvoiceTotal != 0 {
		w.WriteHeader(http.StatusForbidden)
		_ = json.NewEncoder(w).Encode(map[string]any{"message": "No se puede eliminar la factura"})
		return
	}

	if err = rr.DB.DeleteInvoice(parsedId); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (rf *Router) UpdateInvoice(w http.ResponseWriter, r *http.Request) {
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
			_ = json.NewEncoder(w).Encode(map[string]any{"message": "Proyecto no encontrado"})
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(err)
		return
	}

	invoice := types.InvoiceUpdate{
		Id:            parsedId,
		SupplierId:    inv.Supplier.Id,
		ProjectId:     inv.Project.Id,
		InvoiceNumber: inv.InvoiceNumber,
		InvoiceDate:   inv.InvoiceDate,
	}

	p := make(map[string]any)
	if err = json.NewDecoder(r.Body).Decode(&p); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(err)
		return
	}

	if val, ok := p["invoice_number"].(string); ok {
		invoice.InvoiceNumber = val
	}

	if val, ok := p["invoice_date"].(string); ok {
		invoice.InvoiceDate, err = time.ParseInLocation("2006-01-02", val, time.Local)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]any{"message": "La fecha de la factura es inválida"})
			return
		}
	}

	if err = rf.DB.UpdateInvoice(invoice); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
