package router

import (
	"encoding/json"
	"net/http"

	"github.com/alcb1310/final-bca-go/internal/types"
)

func (rf *Router) GetItems(w http.ResponseWriter, r *http.Request) {
	items := []types.Items{}
	var err error

	if items, err = rf.DB.GetItems(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(items)
}

func (rf *Router) CreateItem(w http.ResponseWriter, r *http.Request) {
	errorResponse := make(map[string]any)
	p := make(map[string]any)
	var item types.Items
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

	if item.Code, ok = p["code"].(string); !ok {
		errorResponse["code"] = "El código es obligatorio"
	} else if len(item.Code) == 0 {
		errorResponse["code"] = "El código es obligatorio"
	}

	if item.Name, ok = p["name"].(string); !ok {
		errorResponse["name"] = "El nombre es obligatorio"
	} else if len(item.Name) == 0 {
		errorResponse["name"] = "El nombre es obligatorio"
	}

	if item.Unit, ok = p["unit"].(string); !ok {
		errorResponse["unit"] = "La unidad es obligatoria"
	} else if len(item.Unit) == 0 {
		errorResponse["unit"] = "La unidad es obligatoria"
	}

	if len(errorResponse) > 0 {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errorResponse)
		return
	}

	w.WriteHeader(http.StatusNotImplemented)
	_ = json.NewEncoder(w).Encode(item)
}
