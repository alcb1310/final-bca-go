package router

import (
	"encoding/json"
	"net/http"
)

func (rf *Router) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]any{"message": "OK"})
}
