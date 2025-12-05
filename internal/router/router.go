package router

import (
	"github.com/go-chi/chi/v5"
)
type Router struct {
	Router *chi.Mux
}

func NewRouter() *Router {
	r := chi.NewRouter()

	return &Router{
		Router: r,
	}
}
