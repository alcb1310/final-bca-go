package router

import (
	"encoding/json"
	"net/http"

	"github.com/alcb1310/final-bca-go/internal/types"
)

func (rf *Router) GetSuppliers(w http.ResponseWriter, r *http.Request) {
	s, err := rf.DB.GetSuppliers()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(s)
}

func (rf *Router) CreateSupplier(w http.ResponseWriter, r *http.Request) {
	errorResponse := make(map[string]any)
	p := make(map[string]any)
	var supplier types.Supplier
	var err error
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

	if supplier.Name, ok = p["name"].(string); !ok {
		errorResponse["name"] = "El nombre es obligatorio"
	} else if len(supplier.Name) == 0 {
		errorResponse["name"] = "El nombre es obligatorio"
	}

	if supplier.SupplierId, ok = p["supplier_id"].(string); !ok {
		errorResponse["supplier_id"] = "El RUC es obligatorio"
	} else if len(supplier.SupplierId) == 0 {
		errorResponse["supplier_id"] = "El RUC es obligatorio"
	}

	var val string

	if val, ok = p["contact_name"].(string); !ok {
		supplier.ContactName.Valid = false
	} else {
		if len(val) == 0 {
			supplier.ContactName.Valid = false
		} else {
			supplier.ContactName.Valid = true
			supplier.ContactName.String = val
		}
	}

	if val, ok = p["contact_email"].(string); !ok {
		supplier.ContactEmail.Valid = false
	} else {
		if len(val) == 0 {
			supplier.ContactName.Valid = false
		} else {
			supplier.ContactEmail.Valid = true
			supplier.ContactEmail.String = val
		}
	}

	if val, ok = p["contact_phone"].(string); !ok {
		supplier.ContactPhone.Valid = false
	} else {
		if len(val) == 0 {
			supplier.ContactName.Valid = false
		} else {
			supplier.ContactPhone.Valid = true
			supplier.ContactPhone.String = val
		}
	}

	if len(errorResponse) > 0 {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errorResponse)
		return
	}
	w.WriteHeader(http.StatusNotImplemented)
}
