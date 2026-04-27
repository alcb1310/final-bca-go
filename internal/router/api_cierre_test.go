package router_test

import (
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/alcb1310/final-bca-go/internal/router"
	"github.com/alcb1310/final-bca-go/internal/types"
	"github.com/alcb1310/final-bca-go/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestApiCreateClosure(t *testing.T) {
	db := mocks.NewService(t)
	s := router.Router{DB: db}
	assert.NotNil(t, s)
	s.Router()

	server := httptest.NewServer(s.Router())
	defer server.Close()

	testURL := "/api/v2/cierre"
	projectId := uuid.New()
	dateString := "2026-01-31"
	date, err := time.ParseInLocation("2006-01-02", dateString, time.Local)
	assert.NoError(t, err)

	testData := []struct {
		name    string
		form    map[string]any
		status  int
		body    map[string]any
		project *mocks.Service_GetProject_Call
		closure *mocks.Service_GenerateClosure_Call
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
			name:   "should pass a project",
			form:   map[string]any{},
			status: http.StatusBadRequest,
			body: map[string]any{
				"project_id": "El proyecto es obligatorio",
				"date":       "La fecha es obligatoria",
			},
		},
		{
			name: "should pass a valid project id",
			form: map[string]any{
				"project_id": "invalid",
			},
			status: http.StatusBadRequest,
			body: map[string]any{
				"project_id": "El código del proyecto es inválido",
				"date":       "La fecha es obligatoria",
			},
		},
		{
			name: "should pass an existing project",
			form: map[string]any{
				"project_id": projectId.String(),
			},
			status: http.StatusNotFound,
			body: map[string]any{
				"project_id": "El proyecto no existe",
				"date":       "La fecha es obligatoria",
			},
			project: db.EXPECT().GetProject(projectId).Return(types.Project{}, sql.ErrNoRows),
		},
		{
			name: "should pass a valid date",
			form: map[string]any{
				"project_id": projectId.String(),
				"date":       "invalid",
			},
			status: http.StatusBadRequest,
			body: map[string]any{
				"date": "La fecha es inválida",
			},
			project: db.EXPECT().GetProject(projectId).Return(types.Project{}, nil),
		},
		{
			name: "should pass a valid month",
			form: map[string]any{
				"project_id": projectId.String(),
				"date":       "2026-13-13",
			},
			status: http.StatusBadRequest,
			body: map[string]any{
				"date": "La fecha es inválida",
			},
			project: db.EXPECT().GetProject(projectId).Return(types.Project{}, nil),
		},
		{
			name: "should pass a valid day",
			form: map[string]any{
				"project_id": projectId.String(),
				"date":       "2026-12-33",
			},
			status: http.StatusBadRequest,
			body: map[string]any{
				"date": "La fecha es inválida",
			},
			project: db.EXPECT().GetProject(projectId).Return(types.Project{}, nil),
		},
		{
			name: "should generate the closure",
			form: map[string]any{
				"project_id": projectId.String(),
				"date":       dateString,
			},
			status:  http.StatusNoContent,
			body:    map[string]any{},
			project: db.EXPECT().GetProject(projectId).Return(types.Project{}, nil),
			closure: db.EXPECT().GenerateClosure(projectId, date).Return(nil),
		},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			var read io.Reader = nil

			if tt.project != nil {
				tt.project.Times(1)
			}

			if tt.closure != nil {
				tt.closure.Times(1)
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
			s.Router().ServeHTTP(res, req)
			assert.Equal(t, tt.status, res.Code)

			body, err := io.ReadAll(res.Body)
			assert.NoError(t, err)

			if res.Code == http.StatusNoContent {
				return
			}

			mapBody := make(map[string]any)
			err = json.Unmarshal(body, &mapBody)
			assert.NoError(t, err)

			for k, v := range tt.body {
				assert.Equal(t, v, mapBody[k])
			}
		})
	}
}
