package router

import (
	"encoding/json"
	"net/http"

	"github.com/alcb1310/final-bca-go/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Router struct {
	Router *chi.Mux
	DB     database.Service
}

func NewRouter(db database.Service) *Router {

	r := chi.NewRouter()

	return &Router{
		Router: r,
		DB:     db,
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

		r.Route("/api/v2", func(r chi.Router) {
			r.Get("/health", rf.HealthCheck)

			r.Route("/projects", func(r chi.Router) {
				r.Get("/", rf.GetProjects)
				r.Post("/", rf.CreateProject)
				r.Put("/{id}", rf.UpdateProject)
			})

			r.Route("/suppliers", func(r chi.Router) {
				r.Get("/", rf.GetSuppliers)
				r.Post("/", rf.CreateSupplier)
				r.Put("/{id}", rf.UpdateSupplier)
			})

			r.Route("/budget-items", func(r chi.Router) {
				r.Get("/", rf.GetBudgetItems)
				r.Post("/", rf.CreateBudgetItem)
				r.Put("/{id}", rf.UpdateBudgetItem)
			})

			r.Route("/categories", func(r chi.Router) {
				r.Get("/", rf.GetCategories)
				r.Post("/", rf.CreateCategory)
				r.Put("/{id}", rf.UpdateCategory)
			})

			r.Route("/materials", func(r chi.Router) {
				r.Get("/", rf.GetMaterials)
				r.Post("/", rf.CreateMaterial)
				r.Put("/{id}", rf.UpdateMaterial)
				r.Delete("/{id}", rf.DeleteMaterial)
			})

			r.Route("/items", func(r chi.Router) {
				r.Get("/", rf.GetItems)
				r.Post("/", rf.CreateItem)

				r.Route("/{id}", func(r chi.Router) {
					r.Put("/", rf.UpdateItem)
				})
			})
		})
	})
}
