package router

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/alcb1310/final-bca-go/internal/types"
	"github.com/google/uuid"
)

func (rf *Router) GetBudgets(w http.ResponseWriter, r *http.Request) {
	var err error
	budgets := []types.Budget{}

	if budgets, err = rf.DB.GetBudgets(); err != nil {
		slog.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(err)
		return
	}

	w.WriteHeader(http.StatusNotImplemented)
	_ = json.NewEncoder(w).Encode(budgets)
}

func (rf *Router) CreateBudget(w http.ResponseWriter, r *http.Request) {
	errorResponse := make(map[string]any)
	p := make(map[string]any)
	var budget types.CreateBudget
	var err error
	var val string
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

	if val, ok = p["project_id"].(string); !ok {
		errorResponse["project_id"] = "El proyecto es obligatorio"
	} else {
		budget.ProjectId, err = uuid.Parse(val)
		if err != nil {
			errorResponse["project_id"] = "El código del proyecto es inválido"
		} else {
			if _, err = rf.DB.GetProject(budget.ProjectId); err != nil {
				errorResponse["project_id"] = "El proyecto no existe"
				w.WriteHeader(http.StatusNotFound)
				_ = json.NewEncoder(w).Encode(errorResponse)
				return
			}
		}
	}

	if val, ok = p["budget_item_id"].(string); !ok {
		errorResponse["budget_item_id"] = "La partida es obligatoria"
	} else {
		budget.BudgetItemId, err = uuid.Parse(val)
		if err != nil {
			errorResponse["budget_item_id"] = "El código de la partida es inválido"
		} else {
			if _, err = rf.DB.GetBudgetItem(budget.BudgetItemId); err != nil {
				errorResponse["budget_item_id"] = "La partida no existe"
				w.WriteHeader(http.StatusNotFound)
				_ = json.NewEncoder(w).Encode(errorResponse)
				return
			}
		}
	}

	if budget.Quantity, ok = p["quantity"].(float64); !ok {
		errorResponse["quantity"] = "La cantidad es obligatoria"
	}

	if budget.Cost, ok = p["cost"].(float64); !ok {
		errorResponse["cost"] = "El costo es obligatorio"
	}

	if len(errorResponse) > 0 {
		w.WriteHeader(http.StatusUnprocessableEntity)
		_ = json.NewEncoder(w).Encode(errorResponse)
		return
	}

	w.WriteHeader(http.StatusNotImplemented)
	_ = json.NewEncoder(w).Encode(map[string]any{"message": "Not implemented", "budget": budget})
}
