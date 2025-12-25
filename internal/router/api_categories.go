package router

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/alcb1310/final-bca-go/internal/types"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
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

func (rf *Router) CreateCategory(w http.ResponseWriter, r *http.Request) {
	errorResponse := make(map[string]any)
	p := make(map[string]any)
	var category types.Category
	var err error
	var ok bool

	if r.Body == http.NoBody || r.Body == nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		errorResponse["message"] = "Falta el cuerpo de la solicitud"
		_ = json.NewEncoder(w).Encode(errorResponse)
		return
	}

	if err = json.NewDecoder(r.Body).Decode(&p); err != nil {
		errorResponse["message"] = "Cuerpo de la solicitud no válido"
	}

	if category.Name, ok = p["name"].(string); !ok {
		errorResponse["name"] = "El nombre es obligatorio"
	} else if len(category.Name) == 0 {
		errorResponse["name"] = "El nombre es obligatorio"
	}

	if len(errorResponse) > 0 {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errorResponse)
		return
	}

	if err = rf.DB.CreateCategory(category); err != nil {
		var e *pgconn.PgError
		if errors.As(err, &e) && e.Code == "23505" {
			w.WriteHeader(http.StatusConflict)
			_ = json.NewEncoder(w).Encode(map[string]any{"message": "La categoría ya existe"})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]any{"message": "Categoría creada"})
}

func (rf *Router) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	pId := chi.URLParam(r, "id")
	parsedId, err := uuid.Parse(pId)
	if err != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		_ = json.NewEncoder(w).Encode(map[string]any{"message": "Id inválido"})
		return
	}

	category, err := rf.DB.GetCategory(parsedId)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]any{"message": "Categoría no encontrada"})
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]any{"message": "Error al buscar la categoría", "err": err})
		return
	}

	p := make(map[string]any)

	if err = json.NewDecoder(r.Body).Decode(&p); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]any{"message": "Cuerpo de la solicitud no válido"})
		return
	}
	valStr, ok := p["name"].(string)
	if ok {
		category.Name = valStr
	}

	if err := rf.DB.UpdateCategory(category); err != nil {
		var e *pgconn.PgError
		if errors.As(err, &e) && e.Code == "23505" {
			w.WriteHeader(http.StatusConflict)
			_ = json.NewEncoder(w).Encode(map[string]any{"message": "La categoría ya existe"})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(err)
		return
	}

	w.WriteHeader(http.StatusNotImplemented)
	_ = json.NewEncoder(w).Encode(map[string]any{"message": "Not implemented", "category": category})
}
