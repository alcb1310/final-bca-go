package tests

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
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

func TestApiInvoiceDetails(t *testing.T) {
	ctx := context.Background()
	pgContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithOrderedInitScripts(
			filepath.Join("..", "schema", "tables.sql"),
			filepath.Join("scripts", "seed_projects.sql"),
			filepath.Join("scripts", "seed_suppliers.sql"),
			filepath.Join("scripts", "seed_budget-items.sql"),
			filepath.Join("scripts", "seed_budget.sql"),
			filepath.Join("scripts", "seed_invoices.sql"),
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
		slog.Error("TestApiInvoiceDetails, failed to run pgContainer", "error", err)
		panic(err)
	}

	t.Cleanup(func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			t.Fatalf("TestApiInvoiceDetails, failed to terminate pgContainer: %v", err)
		}
	})

	s, err := createServer(t, ctx, pgContainer)
	assert.NoError(t, err)
	if err != nil {
		slog.Error("TestApiInvoiceDetails, failed to create server", "error", err)
		panic(err)
	}
	s.GenerateRoutes()
	invoiceId := uuid.MustParse("c3be2956-1c3c-46f7-af14-d28420116f14")
	testURL := fmt.Sprintf("/api/v2/invoices/%s/details", invoiceId)

	t.Run("should have no invoice details", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, testURL, nil)
		assert.NoError(t, err)
		resp := httptest.NewRecorder()
		s.Router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Equal(t, "[]", strings.TrimSpace(resp.Body.String()))
	})

	t.Run("should be able to create an invoice detail", func(t *testing.T) {
		budgetResponse := []types.Budget{
			{
				Project:           types.Project{Id: uuid.MustParse("2118e27d-1ae5-4554-b0ba-2503917a31aa")},
				BudgetItem:        types.BudgetItem{Id: uuid.MustParse("420f8bb3-bc8e-4564-be99-75cd7c1a6ff8")},
				InitialQuantity:   sql.NullFloat64{Float64: 0, Valid: false},
				InitialCost:       sql.NullFloat64{Float64: 0, Valid: false},
				InitialTotal:      4567.5,
				SpentQuantity:     sql.NullFloat64{Float64: 0, Valid: false},
				SpentTotal:        0,
				RemainingQuantity: sql.NullFloat64{Float64: 0, Valid: false},
				RemainingCost:     sql.NullFloat64{Float64: 0, Valid: false},
				RemainingTotal:    4567.5,
				UpdatedBudget:     4567.5,
			},
			{
				Project:           types.Project{Id: uuid.MustParse("2118e27d-1ae5-4554-b0ba-2503917a31aa")},
				BudgetItem:        types.BudgetItem{Id: uuid.MustParse("9abc2426-a92b-46ef-b074-ddbc8ee2df1a")},
				InitialQuantity:   sql.NullFloat64{Float64: 2537.5, Valid: true},
				InitialCost:       sql.NullFloat64{Float64: 1.8, Valid: true},
				InitialTotal:      4567.5,
				SpentQuantity:     sql.NullFloat64{Float64: 0, Valid: true},
				SpentTotal:        0,
				RemainingQuantity: sql.NullFloat64{Float64: 2537.5, Valid: true},
				RemainingCost:     sql.NullFloat64{Float64: 1.8, Valid: true},
				RemainingTotal:    4567.5,
				UpdatedBudget:     4567.5,
			},
			{
				Project:           types.Project{Id: uuid.MustParse("2118e27d-1ae5-4554-b0ba-2503917a31aa")},
				BudgetItem:        types.BudgetItem{Id: uuid.MustParse("439082ad-f1bd-4228-91f2-8e744894ffdc")},
				InitialQuantity:   sql.NullFloat64{Float64: 0, Valid: false},
				InitialCost:       sql.NullFloat64{Float64: 0, Valid: false},
				InitialTotal:      100.0,
				SpentQuantity:     sql.NullFloat64{Float64: 0, Valid: false},
				SpentTotal:        25.0,
				RemainingQuantity: sql.NullFloat64{Float64: 0, Valid: false},
				RemainingCost:     sql.NullFloat64{Float64: 0, Valid: false},
				RemainingTotal:    75.0,
				UpdatedBudget:     100.0,
			},
			{
				Project:           types.Project{Id: uuid.MustParse("2118e27d-1ae5-4554-b0ba-2503917a31aa")},
				BudgetItem:        types.BudgetItem{Id: uuid.MustParse("b4b2e4e4-f22d-402e-9ab5-1d59347cbfcb")},
				InitialQuantity:   sql.NullFloat64{Float64: 4, Valid: true},
				InitialCost:       sql.NullFloat64{Float64: 25, Valid: true},
				InitialTotal:      100.0,
				SpentQuantity:     sql.NullFloat64{Float64: 1, Valid: true},
				SpentTotal:        25.0,
				RemainingQuantity: sql.NullFloat64{Float64: 3, Valid: true},
				RemainingCost:     sql.NullFloat64{Float64: 25, Valid: true},
				RemainingTotal:    75.0,
				UpdatedBudget:     100.0,
			},
			{
				Project:           types.Project{Id: uuid.MustParse("1c6020db-39a0-451d-89ee-fdd20d519828")},
				BudgetItem:        types.BudgetItem{Id: uuid.MustParse("439082ad-f1bd-4228-91f2-8e744894ffdc")},
				InitialQuantity:   sql.NullFloat64{Float64: 0, Valid: false},
				InitialCost:       sql.NullFloat64{Float64: 0, Valid: false},
				InitialTotal:      100.0,
				SpentQuantity:     sql.NullFloat64{Float64: 0, Valid: false},
				SpentTotal:        0.0,
				RemainingQuantity: sql.NullFloat64{Float64: 0, Valid: false},
				RemainingCost:     sql.NullFloat64{Float64: 0, Valid: false},
				RemainingTotal:    100.0,
				UpdatedBudget:     100.0,
			},
			{
				Project:           types.Project{Id: uuid.MustParse("1c6020db-39a0-451d-89ee-fdd20d519828")},
				BudgetItem:        types.BudgetItem{Id: uuid.MustParse("b4b2e4e4-f22d-402e-9ab5-1d59347cbfcb")},
				InitialQuantity:   sql.NullFloat64{Float64: 4, Valid: true},
				InitialCost:       sql.NullFloat64{Float64: 25, Valid: true},
				InitialTotal:      100.0,
				SpentQuantity:     sql.NullFloat64{Float64: 0, Valid: true},
				SpentTotal:        0.0,
				RemainingQuantity: sql.NullFloat64{Float64: 4, Valid: true},
				RemainingCost:     sql.NullFloat64{Float64: 25, Valid: true},
				RemainingTotal:    100.0,
				UpdatedBudget:     100.0,
			},
		}

		form := map[string]any{
			"budget_item_id": "b4b2e4e4-f22d-402e-9ab5-1d59347cbfcb",
			"quantity":       1,
			"cost":           25,
		}

		invoiceDetails := types.InvoiceDetailsResponse{
			Invoice:    types.InvoiceResponse{Id: invoiceId},
			BudgetItem: types.BudgetItem{Id: uuid.MustParse("b4b2e4e4-f22d-402e-9ab5-1d59347cbfcb")},
			Project:    types.Project{Id: uuid.MustParse("2118e27d-1ae5-4554-b0ba-2503917a31aa")},
			Supplier:   types.Supplier{Id: uuid.MustParse("b4b2e4e4-f22d-402e-9ab5-1d59347cbfcb")},
			Quantity:   1,
			Cost:       25,
			Total:      25,
		}

		j, err := json.Marshal(form)
		assert.NoError(t, err)

		req, err := http.NewRequest(http.MethodPost, testURL, strings.NewReader(string(j)))
		assert.NoError(t, err)

		req.Header.Add("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		s.Router.ServeHTTP(resp, req)
		assert.Equal(t, resp.Code, http.StatusCreated)

		savedDetails, err := s.DB.GetInvoiceDetails(invoiceId)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(savedDetails))

		assert.Equal(t, savedDetails[0].BudgetItem.Id, invoiceDetails.BudgetItem.Id)
		assert.Equal(t, savedDetails[0].Quantity, invoiceDetails.Quantity)
		assert.Equal(t, savedDetails[0].Cost, invoiceDetails.Cost)
		assert.Equal(t, savedDetails[0].Total, invoiceDetails.Total)

		budgets, err := s.DB.GetBudgets()
		assert.NoError(t, err)
		assert.Equal(t, 6, len(budgets))

		for i, b := range budgets {
			assert.Equal(t, b.Project.Id, budgetResponse[i].Project.Id)
			assert.Equal(t, b.BudgetItem.Id, budgetResponse[i].BudgetItem.Id)
			assert.Equal(t, b.InitialQuantity, budgetResponse[i].InitialQuantity)
			assert.Equal(t, b.InitialCost, budgetResponse[i].InitialCost)
			assert.Equal(t, b.InitialTotal, budgetResponse[i].InitialTotal)
			assert.Equal(t, b.SpentQuantity, budgetResponse[i].SpentQuantity)
			assert.Equal(t, b.SpentTotal, budgetResponse[i].SpentTotal)
			assert.Equal(t, b.RemainingQuantity, budgetResponse[i].RemainingQuantity)
			assert.Equal(t, b.RemainingCost, budgetResponse[i].RemainingCost)
			assert.Equal(t, b.RemainingTotal, budgetResponse[i].RemainingTotal)
			assert.Equal(t, b.UpdatedBudget, budgetResponse[i].UpdatedBudget)
		}
	})

	t.Run("should create a conflict when creating the same detail twice", func(t *testing.T) {
		form := map[string]any{
			"budget_item_id": "b4b2e4e4-f22d-402e-9ab5-1d59347cbfcb",
			"quantity":       1,
			"cost":           25,
		}

		j, err := json.Marshal(form)
		assert.NoError(t, err)

		req, err := http.NewRequest(http.MethodPost, testURL, strings.NewReader(string(j)))
		assert.NoError(t, err)
		req.Header.Add("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		s.Router.ServeHTTP(resp, req)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusConflict, resp.Code)
		assert.Equal(t, "{\"message\":\"El detalle de la factura ya existe\"}", strings.TrimSpace(resp.Body.String()))
	})

	t.Run("should delete an invoice detail", func(t *testing.T) {
		budgetResponse := []types.Budget{
			{
				Project:           types.Project{Id: uuid.MustParse("2118e27d-1ae5-4554-b0ba-2503917a31aa")},
				BudgetItem:        types.BudgetItem{Id: uuid.MustParse("420f8bb3-bc8e-4564-be99-75cd7c1a6ff8")},
				InitialQuantity:   sql.NullFloat64{Float64: 0, Valid: false},
				InitialCost:       sql.NullFloat64{Float64: 0, Valid: false},
				InitialTotal:      4567.5,
				SpentQuantity:     sql.NullFloat64{Float64: 0, Valid: false},
				SpentTotal:        0,
				RemainingQuantity: sql.NullFloat64{Float64: 0, Valid: false},
				RemainingCost:     sql.NullFloat64{Float64: 0, Valid: false},
				RemainingTotal:    4567.5,
				UpdatedBudget:     4567.5,
			},
			{
				Project:           types.Project{Id: uuid.MustParse("2118e27d-1ae5-4554-b0ba-2503917a31aa")},
				BudgetItem:        types.BudgetItem{Id: uuid.MustParse("9abc2426-a92b-46ef-b074-ddbc8ee2df1a")},
				InitialQuantity:   sql.NullFloat64{Float64: 2537.5, Valid: true},
				InitialCost:       sql.NullFloat64{Float64: 1.8, Valid: true},
				InitialTotal:      4567.5,
				SpentQuantity:     sql.NullFloat64{Float64: 0, Valid: true},
				SpentTotal:        0,
				RemainingQuantity: sql.NullFloat64{Float64: 2537.5, Valid: true},
				RemainingCost:     sql.NullFloat64{Float64: 1.8, Valid: true},
				RemainingTotal:    4567.5,
				UpdatedBudget:     4567.5,
			},
			{
				Project:           types.Project{Id: uuid.MustParse("2118e27d-1ae5-4554-b0ba-2503917a31aa")},
				BudgetItem:        types.BudgetItem{Id: uuid.MustParse("439082ad-f1bd-4228-91f2-8e744894ffdc")},
				InitialQuantity:   sql.NullFloat64{Float64: 0, Valid: false},
				InitialCost:       sql.NullFloat64{Float64: 0, Valid: false},
				InitialTotal:      100.0,
				SpentQuantity:     sql.NullFloat64{Float64: 0, Valid: false},
				SpentTotal:        0.0,
				RemainingQuantity: sql.NullFloat64{Float64: 0, Valid: false},
				RemainingCost:     sql.NullFloat64{Float64: 0, Valid: false},
				RemainingTotal:    100.0,
				UpdatedBudget:     100.0,
			},
			{
				Project:           types.Project{Id: uuid.MustParse("2118e27d-1ae5-4554-b0ba-2503917a31aa")},
				BudgetItem:        types.BudgetItem{Id: uuid.MustParse("b4b2e4e4-f22d-402e-9ab5-1d59347cbfcb")},
				InitialQuantity:   sql.NullFloat64{Float64: 4, Valid: true},
				InitialCost:       sql.NullFloat64{Float64: 25, Valid: true},
				InitialTotal:      100.0,
				SpentQuantity:     sql.NullFloat64{Float64: 0, Valid: true},
				SpentTotal:        0.0,
				RemainingQuantity: sql.NullFloat64{Float64: 4, Valid: true},
				RemainingCost:     sql.NullFloat64{Float64: 25, Valid: true},
				RemainingTotal:    100.0,
				UpdatedBudget:     100.0,
			},
			{
				Project:           types.Project{Id: uuid.MustParse("1c6020db-39a0-451d-89ee-fdd20d519828")},
				BudgetItem:        types.BudgetItem{Id: uuid.MustParse("439082ad-f1bd-4228-91f2-8e744894ffdc")},
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
			{
				Project:           types.Project{Id: uuid.MustParse("1c6020db-39a0-451d-89ee-fdd20d519828")},
				BudgetItem:        types.BudgetItem{Id: uuid.MustParse("b4b2e4e4-f22d-402e-9ab5-1d59347cbfcb")},
				InitialQuantity:   sql.NullFloat64{Float64: 4, Valid: true},
				InitialCost:       sql.NullFloat64{Float64: 25, Valid: true},
				InitialTotal:      100,
				SpentQuantity:     sql.NullFloat64{Float64: 0, Valid: true},
				SpentTotal:        0,
				RemainingQuantity: sql.NullFloat64{Float64: 4, Valid: true},
				RemainingCost:     sql.NullFloat64{Float64: 25, Valid: true},
				RemainingTotal:    100,
				UpdatedBudget:     100,
			},
		}

		form := map[string]any{}
		j, err := json.Marshal(form)
		assert.NoError(t, err)

		req, err := http.NewRequest(
			http.MethodDelete,
			fmt.Sprintf("%s/%s", testURL, "b4b2e4e4-f22d-402e-9ab5-1d59347cbfcb"),
			strings.NewReader(string(j)),
		)
		assert.NoError(t, err)
		req.Header.Add("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		s.Router.ServeHTTP(resp, req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.Code)

		budgets, err := s.DB.GetBudgets()
		assert.NoError(t, err)
		assert.Equal(t, 6, len(budgets))

		for i, b := range budgets {
			assert.Equal(t, b.Project.Id, budgetResponse[i].Project.Id)
			assert.Equal(t, b.BudgetItem.Id, budgetResponse[i].BudgetItem.Id)
			assert.Equal(t, b.InitialQuantity, budgetResponse[i].InitialQuantity)
			assert.Equal(t, b.InitialCost, budgetResponse[i].InitialCost)
			assert.Equal(t, b.InitialTotal, budgetResponse[i].InitialTotal)
			assert.Equal(t, b.SpentQuantity, budgetResponse[i].SpentQuantity)
			assert.Equal(t, b.SpentTotal, budgetResponse[i].SpentTotal)
			assert.Equal(t, b.RemainingQuantity, budgetResponse[i].RemainingQuantity)
			assert.Equal(t, b.RemainingCost, budgetResponse[i].RemainingCost)
			assert.Equal(t, b.RemainingTotal, budgetResponse[i].RemainingTotal)
			assert.Equal(t, b.UpdatedBudget, budgetResponse[i].UpdatedBudget)
		}
	})

	t.Run("should update correctly the budget", func(t *testing.T) {
		invoiceDetails := []types.InvoiceDetailsResponse{
			{
				Project:    types.Project{Id: uuid.MustParse("2118e27d-1ae5-4554-b0ba-2503917a31aa")},
				BudgetItem: types.BudgetItem{Id: uuid.MustParse("b4b2e4e4-f22d-402e-9ab5-1d59347cbfcb")},
				Quantity:   1,
				Cost:       30,
				Total:      30,
			},
		}

		budgetResponse := []types.Budget{
			{
				Project:           types.Project{Id: uuid.MustParse("2118e27d-1ae5-4554-b0ba-2503917a31aa")},
				BudgetItem:        types.BudgetItem{Id: uuid.MustParse("420f8bb3-bc8e-4564-be99-75cd7c1a6ff8")},
				InitialQuantity:   sql.NullFloat64{Float64: 0, Valid: false},
				InitialCost:       sql.NullFloat64{Float64: 0, Valid: false},
				InitialTotal:      4567.5,
				SpentQuantity:     sql.NullFloat64{Float64: 0, Valid: false},
				SpentTotal:        0,
				RemainingQuantity: sql.NullFloat64{Float64: 0, Valid: false},
				RemainingCost:     sql.NullFloat64{Float64: 0, Valid: false},
				RemainingTotal:    4567.5,
				UpdatedBudget:     4567.5,
			},
			{
				Project:           types.Project{Id: uuid.MustParse("2118e27d-1ae5-4554-b0ba-2503917a31aa")},
				BudgetItem:        types.BudgetItem{Id: uuid.MustParse("9abc2426-a92b-46ef-b074-ddbc8ee2df1a")},
				InitialQuantity:   sql.NullFloat64{Float64: 2537.5, Valid: true},
				InitialCost:       sql.NullFloat64{Float64: 1.8, Valid: true},
				InitialTotal:      4567.5,
				SpentQuantity:     sql.NullFloat64{Float64: 0, Valid: true},
				SpentTotal:        0,
				RemainingQuantity: sql.NullFloat64{Float64: 2537.5, Valid: true},
				RemainingCost:     sql.NullFloat64{Float64: 1.8, Valid: true},
				RemainingTotal:    4567.5,
				UpdatedBudget:     4567.5,
			},
			{
				Project:           types.Project{Id: uuid.MustParse("2118e27d-1ae5-4554-b0ba-2503917a31aa")},
				BudgetItem:        types.BudgetItem{Id: uuid.MustParse("439082ad-f1bd-4228-91f2-8e744894ffdc")},
				InitialQuantity:   sql.NullFloat64{Float64: 0, Valid: false},
				InitialCost:       sql.NullFloat64{Float64: 0, Valid: false},
				InitialTotal:      100.0,
				SpentQuantity:     sql.NullFloat64{Float64: 0, Valid: false},
				SpentTotal:        30.0,
				RemainingQuantity: sql.NullFloat64{Float64: 0, Valid: false},
				RemainingCost:     sql.NullFloat64{Float64: 0, Valid: false},
				RemainingTotal:    90.0,
				UpdatedBudget:     120.0,
			},
			{
				Project:           types.Project{Id: uuid.MustParse("2118e27d-1ae5-4554-b0ba-2503917a31aa")},
				BudgetItem:        types.BudgetItem{Id: uuid.MustParse("b4b2e4e4-f22d-402e-9ab5-1d59347cbfcb")},
				InitialQuantity:   sql.NullFloat64{Float64: 4, Valid: true},
				InitialCost:       sql.NullFloat64{Float64: 25, Valid: true},
				InitialTotal:      100.0,
				SpentQuantity:     sql.NullFloat64{Float64: 1, Valid: true},
				SpentTotal:        30.0,
				RemainingQuantity: sql.NullFloat64{Float64: 3, Valid: true},
				RemainingCost:     sql.NullFloat64{Float64: 30, Valid: true},
				RemainingTotal:    90.0,
				UpdatedBudget:     120.0,
			},
			{
				Project:           types.Project{Id: uuid.MustParse("1c6020db-39a0-451d-89ee-fdd20d519828")},
				BudgetItem:        types.BudgetItem{Id: uuid.MustParse("439082ad-f1bd-4228-91f2-8e744894ffdc")},
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
			{
				Project:           types.Project{Id: uuid.MustParse("1c6020db-39a0-451d-89ee-fdd20d519828")},
				BudgetItem:        types.BudgetItem{Id: uuid.MustParse("b4b2e4e4-f22d-402e-9ab5-1d59347cbfcb")},
				InitialQuantity:   sql.NullFloat64{Float64: 4, Valid: true},
				InitialCost:       sql.NullFloat64{Float64: 25, Valid: true},
				InitialTotal:      100,
				SpentQuantity:     sql.NullFloat64{Float64: 0, Valid: true},
				SpentTotal:        0,
				RemainingQuantity: sql.NullFloat64{Float64: 4, Valid: true},
				RemainingCost:     sql.NullFloat64{Float64: 25, Valid: true},
				RemainingTotal:    100,
				UpdatedBudget:     100,
			},
		}

		form := map[string]any{
			"budget_item_id": "b4b2e4e4-f22d-402e-9ab5-1d59347cbfcb",
			"quantity":       1,
			"cost":           30,
		}

		j, err := json.Marshal(form)
		assert.NoError(t, err)

		req, err := http.NewRequest(
			http.MethodPost,
			testURL,
			strings.NewReader(string(j)),
		)
		assert.NoError(t, err)
		req.Header.Add("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		s.Router.ServeHTTP(resp, req)
		assert.Equal(t, resp.Code, http.StatusCreated)

		savedDetails, err := s.DB.GetInvoiceDetail(invoiceId, uuid.MustParse("b4b2e4e4-f22d-402e-9ab5-1d59347cbfcb"))
		assert.NoError(t, err)

		assert.Equal(t, savedDetails.BudgetItem.Id, invoiceDetails[0].BudgetItem.Id)
		assert.Equal(t, savedDetails.Quantity, invoiceDetails[0].Quantity)
		assert.Equal(t, savedDetails.Cost, invoiceDetails[0].Cost)
		assert.Equal(t, savedDetails.Total, invoiceDetails[0].Total)

		budgets, err := s.DB.GetBudgets()
		assert.NoError(t, err)
		assert.Equal(t, 6, len(budgets))

		for i, b := range budgets {
			assert.Equal(t, b.Project.Id, budgetResponse[i].Project.Id)
			assert.Equal(t, b.BudgetItem.Id, budgetResponse[i].BudgetItem.Id)
			assert.Equal(t, b.InitialQuantity, budgetResponse[i].InitialQuantity)
			assert.Equal(t, b.InitialCost, budgetResponse[i].InitialCost)
			assert.Equal(t, b.InitialTotal, budgetResponse[i].InitialTotal)
			assert.Equal(t, b.SpentQuantity, budgetResponse[i].SpentQuantity)
			assert.Equal(t, b.SpentTotal, budgetResponse[i].SpentTotal)
			assert.Equal(t, b.RemainingQuantity, budgetResponse[i].RemainingQuantity)
			assert.Equal(t, b.RemainingCost, budgetResponse[i].RemainingCost)
			assert.Equal(t, b.RemainingTotal, budgetResponse[i].RemainingTotal)
			assert.Equal(t, b.UpdatedBudget, budgetResponse[i].UpdatedBudget)
		}
	})
}
