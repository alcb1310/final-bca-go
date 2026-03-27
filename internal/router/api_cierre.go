package router

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (rf *Router) Closure(w http.ResponseWriter, r *http.Request) {
	if r.Body == http.NoBody || r.Body == nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		_ = json.NewEncoder(w).Encode(map[string]any{"message": "Falta el cuerpo de la solicitud"})
		return
	}

	var b map[string]any
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]any{"message": "Cuerpo de la solicitud no válido"})
		return
	}

	var projectId uuid.UUID
	var err error
	var errorResponse = make(map[string]any)
	var closureDate time.Time
	var ok bool
	var val string

	if val, ok = b["project_id"].(string); !ok {
		errorResponse["project_id"] = "El proyecto es obligatorio"
	} else {
		projectId, err = uuid.Parse(val)
		if err != nil {
			errorResponse["project_id"] = "El código del proyecto es inválido"
		} else {
			if _, err = rf.DB.GetProject(projectId); err != nil {
				w.WriteHeader(http.StatusNotFound)
				errorResponse["message"] = "El proyecto no existe"
				_ = json.NewEncoder(w).Encode(errorResponse)
				return
			}
		}
	}

	if val, ok = b["date"].(string); !ok {
		errorResponse["date"] = "La fecha es obligatoria"
	} else {
		closureDate, err = time.ParseInLocation("2006-01-02", val, time.Local)
		if err != nil {
			errorResponse["date"] = "La fecha es inválida"
		}
	}

	if len(errorResponse) > 0 {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errorResponse)
		return
	}

	if err = rf.DB.GenerateClosure(projectId, closureDate); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
