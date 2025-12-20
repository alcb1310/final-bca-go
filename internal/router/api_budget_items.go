package router

import (
	"encoding/json"
	"net/http"

	"github.com/alcb1310/final-bca-go/internal/types"
)

func (rf *Router) GetBudgetItems(w http.ResponseWriter, r *http.Request) {
	var bi []types.BudgetItem
	var err error

	if bi, err = rf.DB.GetBudgetItems(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		res := map[string]any{"message": err.Error()}
		_ = json.NewEncoder(w).Encode(res)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(bi)
}
