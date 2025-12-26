package router

import (
	"encoding/json"
	"net/http"
)

func (rf *Router) GetMaterials(w http.ResponseWriter, r *http.Request) {
	materials, err := rf.DB.GetMaterials()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(err)
		return
	}
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(materials)
}
