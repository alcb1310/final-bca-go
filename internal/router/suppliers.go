package router

import (
	"encoding/json"
	"net/http"
)

func (rf *Router) GetSuppliers(w http.ResponseWriter, r *http.Request) {
	s, err := rf.DB.GetSuppliers()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(s)
}

func (rf *Router) CreateSupplier(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}
