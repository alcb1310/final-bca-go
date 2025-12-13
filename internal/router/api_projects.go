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

func (rf *Router) GetProjects(w http.ResponseWriter, r *http.Request) {
	var projects []types.Project
	var err error

	if projects, err = rf.DB.GetProjects(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		res := map[string]any{"message": err.Error()}
		_ = json.NewEncoder(w).Encode(res)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(projects)
}

func (rf *Router) CreateProject(w http.ResponseWriter, r *http.Request) {
	errorResponse := make(map[string]any)
	p := make(map[string]any)
	var project types.Project
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

	if project.Name, ok = p["name"].(string); !ok {
		errorResponse["name"] = "El nombre es obligatorio"
	}

	if project.IsActive, ok = p["is_active"].(bool); !ok {
		errorResponse["is_active"] = "El estado del projecto es obligatorio"
	}

	if project.GrossArea, ok = p["gross_area"].(float64); !ok {
		errorResponse["gross_area"] = "El área bruta es obligatorio"
	}

	if project.NetArea, ok = p["net_area"].(float64); !ok {
		errorResponse["net_area"] = "El área neta es obligatorio"
	}

	if len(errorResponse) > 0 {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errorResponse)
		return
	}

	if err = rf.DB.CreateProject(project); err != nil {
		var e *pgconn.PgError
		if errors.As(err, &e) && e.Code == "23505" {
			w.WriteHeader(http.StatusConflict)
			_ = json.NewEncoder(w).Encode(map[string]any{"message": "El projecto ya existe"})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		errorResponse["message"] = err.Error()
		_ = json.NewEncoder(w).Encode(errorResponse)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]any{"message": "Proyecto creado correctamente"})
}

func (rf *Router) UpdateProject(w http.ResponseWriter, r *http.Request) {
	pId := chi.URLParam(r, "id")
	parsedId, err := uuid.Parse(pId)
	if err != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		_ = json.NewEncoder(w).Encode(map[string]any{"message": "Id inválido"})
		return
	}

	project, err := rf.DB.GetProject(parsedId)
	if err != nil {

		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]any{"message": "Proyecto no encontrado"})
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]any{"message": "Error al buscar el projecto", "err": err})
		return
	}

	errorResponse := make(map[string]any)
	p := make(map[string]any)
	if err = json.NewDecoder(r.Body).Decode(&p); err != nil {
		errorResponse["message"] = "Cuerpo de la solicitud no válido"
	}

	valStr, ok := p["name"].(string)
	if ok {
		project.Name = valStr
	}

	valBool, ok := p["is_active"].(bool)
	if ok {
		project.IsActive = valBool
	}

	valFloat, ok := p["gross_area"].(float64)
	if ok {
		project.GrossArea = valFloat
	}

	valFloat, ok = p["net_area"].(float64)
	if ok {
		project.NetArea = valFloat
	}

	if len(errorResponse) == 0 {
		w.WriteHeader(http.StatusUnprocessableEntity)
		_ = json.NewEncoder(w).Encode(errorResponse)
		return
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
	_ = json.NewEncoder(w).Encode(map[string]any{"message": "Método no permitido", "project": project})
}
