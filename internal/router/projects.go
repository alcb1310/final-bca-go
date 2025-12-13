package router

import (
	"encoding/json"
	"net/http"

	"github.com/alcb1310/final-bca-go/internal/types"
)

func (rf *Router) CreateProject(w http.ResponseWriter, r *http.Request) {
	errorResponse := make(map[string]any)
	p := make(map[string]any)
	var project types.Project
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

	if project.Name, ok = p["name"].(string); !ok {
		errorResponse["name"] = "El nombre es obligatorio"
	}

	if project.IsActive, ok = p["is_active"].(bool); !ok {
		errorResponse["is_active"] = "El estado del projecto es obligatorio"
	}

	if project.GrossArea, ok = p["gross_area"].(float64); !ok {
		errorResponse["gross_area"] = "El área bruta es obligatorio"
	}

	if project.NetArea, ok = p["net_area"].(float64); !ok {
		errorResponse["net_area"] = "El área neta es obligatorio"
	}

	if len(errorResponse) > 0 {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errorResponse)
		return
	}

	if err = rf.DB.CreateProject(project); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errorResponse["message"] = err.Error()
		_ = json.NewEncoder(w).Encode(errorResponse)
		return
	}

	w.WriteHeader(http.StatusNotImplemented)
	_ = json.NewEncoder(w).Encode(map[string]any{"message": "No implementado"})
}
