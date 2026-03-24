package router

import (
	"encoding/json"
	"net/http"
)

func (rf *Router) Closure(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	_ = json.NewEncoder(w).Encode(map[string]any{"message": "Not implemented"})
}
