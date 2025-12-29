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

func (rf *Router) GetItems(w http.ResponseWriter, r *http.Request) {
	items := []types.Items{}
	var err error

	if items, err = rf.DB.GetItems(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(items)
}

func (rf *Router) CreateItem(w http.ResponseWriter, r *http.Request) {
	errorResponse := make(map[string]any)
	p := make(map[string]any)
	var item types.Items
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

	if item.Code, ok = p["code"].(string); !ok {
		errorResponse["code"] = "El código es obligatorio"
	} else if len(item.Code) == 0 {
		errorResponse["code"] = "El código es obligatorio"
	}

	if item.Name, ok = p["name"].(string); !ok {
		errorResponse["name"] = "El nombre es obligatorio"
	} else if len(item.Name) == 0 {
		errorResponse["name"] = "El nombre es obligatorio"
	}

	if item.Unit, ok = p["unit"].(string); !ok {
		errorResponse["unit"] = "La unidad es obligatoria"
	} else if len(item.Unit) == 0 {
		errorResponse["unit"] = "La unidad es obligatoria"
	}

	if len(errorResponse) > 0 {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errorResponse)
		return
	}

	if err = rf.DB.CreateItem(item); err != nil {
		var e *pgconn.PgError
		if errors.As(err, &e) && e.Code == "23505" {
			w.WriteHeader(http.StatusConflict)
			_ = json.NewEncoder(w).Encode(map[string]any{"message": "El rubro ya existe"})
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]any{"message": "Rubro creado"})
}

func (rf *Router) UpdateItem(w http.ResponseWriter, r *http.Request) {
	pId := chi.URLParam(r, "id")
	parsedId, err := uuid.Parse(pId)
	if err != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		_ = json.NewEncoder(w).Encode(map[string]any{"message": "Id inválido"})
		return
	}

	item, err := rf.DB.GetItem(parsedId)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]any{"message": "Item no encontrado"})
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(err)
		return
	}

	var val string
	var ok bool
	errorResponse := make(map[string]any)
	p := make(map[string]any)

	if r.Body == http.NoBody || r.Body == nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		errorResponse["message"] = "Falta el cuerpo de la solicitud"
		_ = json.NewEncoder(w).Encode(errorResponse)
		return
	}

	if err = json.NewDecoder(r.Body).Decode(&p); err != nil {
		errorResponse["message"] = "Cuerpo de la solicitud no válido"
	}

	if val, ok = p["code"].(string); ok && len(val) > 0 {
		item.Code = val
	}

	if val, ok = p["name"].(string); ok && len(val) > 0 {
		item.Name = val
	}

	if val, ok = p["unit"].(string); ok && len(val) > 0 {
		item.Unit = val
	}

	if err = rf.DB.UpdateItem(item); err != nil {
		var e *pgconn.PgError
		if errors.As(err, &e) && e.Code == "23505" {
			w.WriteHeader(http.StatusConflict)
			_ = json.NewEncoder(w).Encode(map[string]any{"message": "El rubro ya existe"})
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
