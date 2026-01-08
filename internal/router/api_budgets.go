package router

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/alcb1310/final-bca-go/internal/types"
)

func (rf *Router) GetBudgets(w http.ResponseWriter, r *http.Request) {
	var err error
	budgets := []types.Budget{}

	if budgets, err = rf.DB.GetBudgets(); err != nil {
		slog.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(err)
		return
	}

	w.WriteHeader(http.StatusNotImplemented)
	_ = json.NewEncoder(w).Encode(budgets)
}

func (rf *Router) CreateBudget(w http.ResponseWriter, r *http.Request) {
	errorResponse := make(map[string]any)
	p := make(map[string]any)
	var budget types.CreateBudget
	var err error
	// var ok bool

	if r.Body == http.NoBody || r.Body == nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		errorResponse["message"] = "Falta el cuerpo de la solicitud"
		_ = json.NewEncoder(w).Encode(errorResponse)
		return
	}

	if err = json.NewDecoder(r.Body).Decode(&p); err != nil {
		errorResponse["message"] = "Cuerpo de la solicitud no válido"
	}
	w.WriteHeader(http.StatusNotImplemented)
	_ = json.NewEncoder(w).Encode(map[string]any{"message": "Not implemented", "budget": budget})
}
