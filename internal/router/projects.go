package router

import (
	"encoding/json"
	"net/http"
)

func (rf *Router) CreateProject(w http.ResponseWriter, r *http.Request) {
	errorResponse := make(map[string]any)

	if r.Body == http.NoBody || r.Body == nil {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse["message"] = "Falta el cuerpo de la solicitud"
		_ = json.NewEncoder(w).Encode(errorResponse)
		return
	}

	w.WriteHeader(http.StatusNotImplemented)
	_ = json.NewEncoder(w).Encode(map[string]any{"message": "No implementado"})
}
