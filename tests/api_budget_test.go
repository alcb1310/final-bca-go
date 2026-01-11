package tests

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/alcb1310/final-bca-go/internal/types"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestApiBudget(t *testing.T) {
	result := map[string]types.SaveBudget{
		"500": {
			ProjectId:         uuid.MustParse("2118e27d-1ae5-4554-b0ba-2503917a31aa"),
			BudgetItemId:      uuid.MustParse("439082ad-f1bd-4228-91f2-8e744894ffdc"),
			InitialQuantity:   sql.NullFloat64{Float64: 0, Valid: false},
			InitialCost:       sql.NullFloat64{Float64: 0, Valid: false},
			InitialTotal:      100,
			SpentQuantity:     sql.NullFloat64{Float64: 0, Valid: false},
			SpentTotal:        0,
			RemainingQuantity: sql.NullFloat64{Float64: 0, Valid: false},
			RemainingCost:     sql.NullFloat64{Float64: 0, Valid: false},
			RemainingTotal:    100,
			UpdatedBudget:     100,
		},
		"500.1": {
			ProjectId:         uuid.MustParse("2118e27d-1ae5-4554-b0ba-2503917a31aa"),
			BudgetItemId:      uuid.MustParse("b4b2e4e4-f22d-402e-9ab5-1d59347cbfcb"),
			InitialQuantity:   sql.NullFloat64{Float64: 10, Valid: true},
			InitialCost:       sql.NullFloat64{Float64: 10, Valid: true},
			InitialTotal:      100,
			SpentQuantity:     sql.NullFloat64{Float64: 0, Valid: true},
			SpentTotal:        0,
			RemainingQuantity: sql.NullFloat64{Float64: 10, Valid: true},
			RemainingCost:     sql.NullFloat64{Float64: 10, Valid: true},
			RemainingTotal:    100,
			UpdatedBudget:     100,
		},
		"200": {
			ProjectId:         uuid.MustParse("2118e27d-1ae5-4554-b0ba-2503917a31aa"),
			BudgetItemId:      uuid.MustParse("420f8bb3-bc8e-4564-be99-75cd7c1a6ff8"),
			InitialQuantity:   sql.NullFloat64{Float64: 0, Valid: false},
			InitialCost:       sql.NullFloat64{Float64: 0, Valid: false},
			InitialTotal:      50,
			SpentQuantity:     sql.NullFloat64{Float64: 0, Valid: false},
			SpentTotal:        0,
			RemainingQuantity: sql.NullFloat64{Float64: 0, Valid: false},
			RemainingCost:     sql.NullFloat64{Float64: 0, Valid: false},
			RemainingTotal:    50,
			UpdatedBudget:     50,
		},
		"200.1": {
			ProjectId:         uuid.MustParse("2118e27d-1ae5-4554-b0ba-2503917a31aa"),
			BudgetItemId:      uuid.MustParse("9abc2426-a92b-46ef-b074-ddbc8ee2df1a"),
			InitialQuantity:   sql.NullFloat64{Float64: 20, Valid: true},
			InitialCost:       sql.NullFloat64{Float64: 2.5, Valid: true},
			InitialTotal:      50,
			SpentQuantity:     sql.NullFloat64{Float64: 0, Valid: true},
			SpentTotal:        0,
			RemainingQuantity: sql.NullFloat64{Float64: 20, Valid: true},
			RemainingCost:     sql.NullFloat64{Float64: 2.5, Valid: true},
			RemainingTotal:    50,
			UpdatedBudget:     50,
		},
	}

	ctx := context.Background()
	testUrl := "/api/v2/budgets"
	t.Log(testUrl)

	pgContainer, err := postgres.Run(ctx,
		"postgres:18-alpine",
		postgres.WithOrderedInitScripts(
			filepath.Join("..", "schema", "tables.sql"),
			filepath.Join("scripts", "seed_projects.sql"),
			filepath.Join("scripts", "seed_budget-items.sql"),
		),
		postgres.WithDatabase("testbca"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(15*time.Second)),
	)

	assert.NoError(t, err)
	if err != nil {
		if pgContainer != nil {
			err := pgContainer.Terminate(ctx)
			assert.NoError(t, err)
		}
		return
	}

	t.Cleanup(func() {
		err := pgContainer.Terminate(ctx)
		assert.NoError(t, err)
	})

	s, err := createServer(t, ctx, pgContainer)
	assert.NoError(t, err)
	if err != nil {
		return
	}
	s.GenerateRoutes()

	t.Run("should have no budgets", func(t *testing.T) {
		req, err := http.NewRequest("GET", testUrl, nil)
		assert.NoError(t, err)
		res := httptest.NewRecorder()
		s.Router.ServeHTTP(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
		var r []any
		err = json.Unmarshal(res.Body.Bytes(), &r)
		assert.NoError(t, err)
		assert.Equal(t, 0, len(r))
		assert.Equal(t, "[]", strings.TrimSpace(res.Body.String()))
	})

	t.Run("should be able to create a budget", func(t *testing.T) {
		form := map[string]any{
			"project_id":     "2118e27d-1ae5-4554-b0ba-2503917a31aa",
			"budget_item_id": "b4b2e4e4-f22d-402e-9ab5-1d59347cbfcb",
			"quantity":       10,
			"cost":           10,
		}

		j, err := json.Marshal(form)
		assert.NoError(t, err)
		req, err := http.NewRequest(http.MethodPost, testUrl, strings.NewReader(string(j)))
		assert.NoError(t, err)
		req.Header.Add("Content-Type", "application/json")
		res := httptest.NewRecorder()
		s.Router.ServeHTTP(res, req)

		assert.Equal(t, http.StatusCreated, res.Code)

		body, err := io.ReadAll(res.Body)
		assert.NoError(t, err)
		mapBody := make(map[string]any)
		err = json.Unmarshal(body, &mapBody)
		assert.NoError(t, err)
		assert.Equal(t, "Presupuesto creado", mapBody["message"])

		req, err = http.NewRequest("GET", testUrl, nil)
		assert.NoError(t, err)
		res = httptest.NewRecorder()
		s.Router.ServeHTTP(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
		var budgets []types.Budget
		err = json.Unmarshal(res.Body.Bytes(), &budgets)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(budgets))

		for _, budget := range budgets {
			assert.Equal(t, result[budget.BudgetItem.Code].ProjectId, budget.Project.Id)
			assert.Equal(t, result[budget.BudgetItem.Code].BudgetItemId, budget.BudgetItem.Id)
			assert.Equal(t, result[budget.BudgetItem.Code].InitialQuantity, budget.InitialQuantity)
			assert.Equal(t, result[budget.BudgetItem.Code].InitialCost, budget.InitialCost)
			assert.Equal(t, result[budget.BudgetItem.Code].InitialTotal, budget.InitialTotal)
			assert.Equal(t, result[budget.BudgetItem.Code].SpentQuantity, budget.SpentQuantity)
			assert.Equal(t, result[budget.BudgetItem.Code].SpentTotal, budget.SpentTotal)
			assert.Equal(t, result[budget.BudgetItem.Code].RemainingQuantity, budget.RemainingQuantity)
			assert.Equal(t, result[budget.BudgetItem.Code].RemainingCost, budget.RemainingCost)
			assert.Equal(t, result[budget.BudgetItem.Code].RemainingTotal, budget.RemainingTotal)
			assert.Equal(t, result[budget.BudgetItem.Code].UpdatedBudget, budget.UpdatedBudget)

		}

		form = map[string]any{
			"project_id":     "2118e27d-1ae5-4554-b0ba-2503917a31aa",
			"budget_item_id": "9abc2426-a92b-46ef-b074-ddbc8ee2df1a",
			"quantity":       20,
			"cost":           2.5,
		}
		j, _ = json.Marshal(form)
		req, err = http.NewRequest(http.MethodPost, testUrl, strings.NewReader(string(j)))
		assert.NoError(t, err)
		req.Header.Add("Content-Type", "application/json")
		res = httptest.NewRecorder()
		s.Router.ServeHTTP(res, req)
		assert.Equal(t, http.StatusCreated, res.Code)

		req, err = http.NewRequest("GET", testUrl, nil)
		assert.NoError(t, err)
		res = httptest.NewRecorder()
		s.Router.ServeHTTP(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
		err = json.Unmarshal(res.Body.Bytes(), &budgets)
		assert.NoError(t, err)
		assert.Equal(t, 4, len(budgets))

		for _, budget := range budgets {
			assert.Equal(t, result[budget.BudgetItem.Code].ProjectId, budget.Project.Id)
			assert.Equal(t, result[budget.BudgetItem.Code].BudgetItemId, budget.BudgetItem.Id)
			assert.Equal(t, result[budget.BudgetItem.Code].InitialQuantity, budget.InitialQuantity)
			assert.Equal(t, result[budget.BudgetItem.Code].InitialCost, budget.InitialCost)
			assert.Equal(t, result[budget.BudgetItem.Code].InitialTotal, budget.InitialTotal)
			assert.Equal(t, result[budget.BudgetItem.Code].SpentQuantity, budget.SpentQuantity)
			assert.Equal(t, result[budget.BudgetItem.Code].SpentTotal, budget.SpentTotal)
			assert.Equal(t, result[budget.BudgetItem.Code].RemainingQuantity, budget.RemainingQuantity)
			assert.Equal(t, result[budget.BudgetItem.Code].RemainingCost, budget.RemainingCost)
			assert.Equal(t, result[budget.BudgetItem.Code].RemainingTotal, budget.RemainingTotal)
			assert.Equal(t, result[budget.BudgetItem.Code].UpdatedBudget, budget.UpdatedBudget)

		}
	})

	t.Run("it should create a conflict", func(t *testing.T) {
		form := map[string]any{
			"project_id":     "2118e27d-1ae5-4554-b0ba-2503917a31aa",
			"budget_item_id": "b4b2e4e4-f22d-402e-9ab5-1d59347cbfcb",
			"quantity":       10,
			"cost":           10,
		}
		j, _ := json.Marshal(form)
		req, err := http.NewRequest(http.MethodPost, testUrl, strings.NewReader(string(j)))
		assert.NoError(t, err)
		req.Header.Add("Content-Type", "application/json")
		res := httptest.NewRecorder()
		s.Router.ServeHTTP(res, req)
		assert.Equal(t, http.StatusConflict, res.Code)

		body, err := io.ReadAll(res.Body)
		assert.NoError(t, err)
		mapBody := make(map[string]any)
		err = json.Unmarshal(body, &mapBody)
		assert.NoError(t, err)
		assert.Equal(t, "Ya existe un presupuesto con ese proyecto y partida", mapBody["message"])
	})

}
