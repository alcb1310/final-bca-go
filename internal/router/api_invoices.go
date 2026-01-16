package router

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/alcb1310/final-bca-go/internal/types"
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
				errorResponse["project_id"] = "El proyecto no existe"
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
				errorResponse["supplier_id"] = "El proveedor no existe"
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

	w.WriteHeader(http.StatusNotImplemented)
	_ = json.NewEncoder(w).Encode(invoice)
}
