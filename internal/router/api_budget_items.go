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

func (rf *Router) GetBudgetItems(w http.ResponseWriter, r *http.Request) {
	var bi []types.BudgetItem
	var err error

	if bi, err = rf.DB.GetBudgetItems(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		res := map[string]any{"message": err.Error()}
		_ = json.NewEncoder(w).Encode(res)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(bi)
}

func (rf *Router) CreateBudgetItem(w http.ResponseWriter, r *http.Request) {
	errorResponse := make(map[string]any)
	p := make(map[string]any)
	var budgetItem types.CreateBudgetItem
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

	if budgetItem.Code, ok = p["code"].(string); !ok {
		errorResponse["code"] = "El código es obligatorio"
	} else if len(budgetItem.Code) == 0 {
		errorResponse["code"] = "El códigio es obligatorio"
	}

	if budgetItem.Name, ok = p["name"].(string); !ok {
		errorResponse["name"] = "El nombre es obligatorio"
	} else if len(budgetItem.Name) == 0 {
		errorResponse["name"] = "El nombre es obligatorio"
	}

	if budgetItem.Accumulate, ok = p["accumulate"].(bool); !ok {
		errorResponse["accumulate"] = "Debe indicar si acumula o no"
	}

	if val, ok := p["parent_id"].(string); ok {
		uuidVal, err := uuid.Parse(val)
		if err != nil {
			errorResponse["parent_id"] = "Id inválido"
		} else {
			budgetItem.ParentId.UUID = uuidVal
			budgetItem.ParentId.Valid = true
		}
	} else {
		budgetItem.ParentId.Valid = false
	}

	if len(errorResponse) > 0 {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errorResponse)
		return
	}

	if err := rf.DB.CreateBudgetItem(budgetItem); err != nil {
		var e *pgconn.PgError
		if errors.As(err, &e) && e.Code == "23505" {
			w.WriteHeader(http.StatusConflict)
			_ = json.NewEncoder(w).Encode(map[string]any{"message": "La partida ya existe"})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		errorResponse["message"] = err.Error()
		_ = json.NewEncoder(w).Encode(errorResponse)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]any{"message": "Partida creada correctamente"})
}

func (rf *Router) UpdateBudgetItem(w http.ResponseWriter, r *http.Request) {
	pId := chi.URLParam(r, "id")
	parsedId, err := uuid.Parse(pId)
	if err != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		_ = json.NewEncoder(w).Encode(map[string]any{"message": "Id inválido"})
		return
	}

	bi, err := rf.DB.GetBudgetItem(parsedId)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]any{"message": "Partida no encontrada"})
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]any{"message": "Error al buscar la partida", "err": err})
		return
	}

	errorResponse := make(map[string]any)
	p := make(map[string]any)
	var budgetItem types.CreateBudgetItem
	var ok, boolVal bool
	var val string

	if err = json.NewDecoder(r.Body).Decode(&p); err != nil {
		errorResponse["message"] = "Cuerpo de la solicitud no válido"
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errorResponse)
		return
	}

	if val, ok = p["code"].(string); ok {
		if len(budgetItem.Code) != 0 {
			budgetItem.Code = val
		}
	}

	if val, ok = p["name"].(string); !ok {
		if len(val) != 0 {
			budgetItem.Name = val
		}
	}

	if boolVal, ok = p["accumulate"].(bool); ok {
		budgetItem.Accumulate = boolVal
	}

	w.WriteHeader(http.StatusNotImplemented)
	_ = json.NewEncoder(w).Encode(map[string]any{"message": "Not implemented", "budget_item": bi})
}
