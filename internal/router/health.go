package router

import (
	"encoding/json"
	"net/http"
)

func (rf *Router) HealthCheck(w http.ResponseWriter, r *http.Request) {
	if !rf.DB.GetHealth() {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]any{"message": "Unable to connect to the database, please check your logs"})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]any{"message": "The database is healthy"})
}
