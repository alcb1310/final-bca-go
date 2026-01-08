package router_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/alcb1310/final-bca-go/internal/router"
	"github.com/alcb1310/final-bca-go/internal/types"
	"github.com/alcb1310/final-bca-go/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestApiCreateBudget(t *testing.T) {
	projectId := uuid.New()
	budgetItemId := uuid.New()
	db := mocks.NewService(t)
	s := router.NewRouter(db)
	s.GenerateRoutes()
	testURL := "/api/v2/budgets"
	testData := []struct {
		name       string
		form       map[string]any
		status     int
		body       map[string]any
		project    *mocks.Service_GetProject_Call
		budgetItem *mocks.Service_GetBudgetItem_Call
	}{
		{
			name:   "should pass a form",
			form:   nil,
			status: http.StatusUnprocessableEntity,
			body: map[string]any{
				"message": "Falta el cuerpo de la solicitud",
			},
		},
		{
			name:   "should pass a body",
			form:   map[string]any{},
			status: http.StatusBadRequest,
			body: map[string]any{
				"project_id":     "El proyecto es obligatorio",
				"budget_item_id": "La partida es obligatoria",
				"quantity":       "La cantidad es obligatoria",
				"cost":           "El costo es obligatorio",
			},
		},
		{
			name: "should pass a non empty project id",
			form: map[string]any{
				"project_id": "",
			},
			status: http.StatusBadRequest,
			body: map[string]any{
				"project_id":     "El código del proyecto es inválido",
				"budget_item_id": "La partida es obligatoria",
				"quantity":       "La cantidad es obligatoria",
				"cost":           "El costo es obligatorio",
			},
		},
		{
			name: "should pass a valid project id",
			form: map[string]any{
				"project_id": "invaldid",
			},
			status: http.StatusBadRequest,
			body: map[string]any{
				"project_id":     "El código del proyecto es inválido",
				"budget_item_id": "La partida es obligatoria",
				"quantity":       "La cantidad es obligatoria",
				"cost":           "El costo es obligatorio",
			},
		},
		{
			name: "should pass a budget item id",
			form: map[string]any{
				"project_id": projectId.String(),
			},
			status:  http.StatusBadRequest,
			project: db.EXPECT().GetProject(projectId).Return(types.Project{}, nil),
			body: map[string]any{
				"budget_item_id": "La partida es obligatoria",
				"quantity":       "La cantidad es obligatoria",
				"cost":           "El costo es obligatorio",
			},
		},
		{
			name: "should pass a non empty budget item id",
			form: map[string]any{
				"project_id":     projectId.String(),
				"budget_item_id": "",
			},
			status:  http.StatusBadRequest,
			project: db.EXPECT().GetProject(projectId).Return(types.Project{}, nil),
			body: map[string]any{
				"budget_item_id": "El código de la partida es inválido",
				"quantity":       "La cantidad es obligatoria",
				"cost":           "El costo es obligatorio",
			},
		},
		{
			name: "should pass a valid budget item id",
			form: map[string]any{
				"project_id":     projectId.String(),
				"budget_item_id": "invalid",
			},
			status:  http.StatusBadRequest,
			project: db.EXPECT().GetProject(projectId).Return(types.Project{}, nil),
			body: map[string]any{
				"budget_item_id": "El código de la partida es inválido",
				"quantity":       "La cantidad es obligatoria",
				"cost":           "El costo es obligatorio",
			},
		},
		{
			name: "should pass a quantity",
			form: map[string]any{
				"project_id":     projectId.String(),
				"budget_item_id": budgetItemId.String(),
			},
			status:     http.StatusBadRequest,
			project:    db.EXPECT().GetProject(projectId).Return(types.Project{}, nil),
			budgetItem: db.EXPECT().GetBudgetItem(budgetItemId).Return(types.BudgetItem{}, nil),
			body: map[string]any{
				"quantity": "La cantidad es obligatoria",
				"cost":     "El costo es obligatorio",
			},
		},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			var read io.Reader = nil

			if tt.budgetItem != nil {
				tt.budgetItem.Times(1)
			}
			if tt.project != nil {
				tt.project.Times(1)
			}

			if tt.form != nil {
				j, err := json.Marshal(tt.form)
				assert.NoError(t, err)
				read = strings.NewReader(string(j))
			}

			req, err := http.NewRequest(http.MethodPost, testURL, read)
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			res := httptest.NewRecorder()
			s.Router.ServeHTTP(res, req)
			assert.Equal(t, tt.status, res.Code)

			body, err := io.ReadAll(res.Body)
			assert.NoError(t, err)

			mapBody := make(map[string]any)
			err = json.Unmarshal(body, &mapBody)
			assert.NoError(t, err)

			assert.Equal(t, len(tt.body), len(mapBody))

			for k, v := range tt.body {
				assert.Equal(t, v, mapBody[k])
			}
		})
	}
}
