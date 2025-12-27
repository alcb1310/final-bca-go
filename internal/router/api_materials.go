package router

import (
	"encoding/json"
	"net/http"

	"github.com/alcb1310/final-bca-go/internal/types"
	"github.com/google/uuid"
)

func (rf *Router) GetMaterials(w http.ResponseWriter, r *http.Request) {
	materials, err := rf.DB.GetMaterials()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(err)
		return
	}
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(materials)
}

func (rf *Router) CreateMaterial(w http.ResponseWriter, r *http.Request) {
	errorResponse := make(map[string]any)
	p := make(map[string]any)
	var material types.Materials
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

	if material.Code, ok = p["code"].(string); !ok {
		errorResponse["code"] = "El código es obligatorio"
	} else if len(material.Code) == 0 {
		errorResponse["code"] = "El código es obligatorio"
	}

	if material.Name, ok = p["name"].(string); !ok {
		errorResponse["name"] = "El nombre es obligatorio"
	} else if len(material.Name) == 0 {
		errorResponse["name"] = "El nombre es obligatorio"
	}

	if material.Unit, ok = p["unit"].(string); !ok {
		errorResponse["unit"] = "La unidad es obligatoria"
	} else if len(material.Unit) == 0 {
		errorResponse["unit"] = "La unidad es obligatoria"
	}

	val := ""
	if val, ok = p["category_id"].(string); !ok {
		errorResponse["category_id"] = "La categoría es obligatoria"
	} else if len(val) == 0 {
		errorResponse["category_id"] = "La categoría es obligatoria"
	} else {
		material.Category.Id, err = uuid.Parse(val)
		if err != nil {
			errorResponse["category_id"] = "La categoría es inválida"
		}
	}

	if len(errorResponse) > 0 {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errorResponse)
		return
	}

	w.WriteHeader(http.StatusNotImplemented)
}
