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
	w.WriteHeader(http.StatusNotImplemented)
}
