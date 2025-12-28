package router

import "net/http"

func (rf *Router) GetItems(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}
