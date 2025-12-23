package router

import (
	"encoding/json"
	"net/http"

	"github.com/alcb1310/final-bca-go/internal/types"
)

func (rf *Router) GetCategories(w http.ResponseWriter, r *http.Request) {
	var categories []types.Category
	var err error

	categories, err = rf.DB.GetCategories()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(categories)
}
