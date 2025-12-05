package router

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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
	rf.Router.Use(middleware.RequestID)
	rf.Router.Use(middleware.Logger)
	rf.Router.Use(middleware.Recoverer)
	rf.Router.Use(middleware.AllowContentType("application/json"))
	rf.Router.Use(contentTypeMiddleware)

	rf.Router.Route("/", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(map[string]any{"message": "Hello World"})
		})
	})
}
