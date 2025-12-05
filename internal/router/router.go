package router

import (
	"net/http"

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

func (rf *Router) GenerateRoutes() {
	rf.Router.Route("/", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hello World"))
		})
	})
}
