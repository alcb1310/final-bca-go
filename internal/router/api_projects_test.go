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
	"github.com/stretchr/testify/assert"
)

func TestApiCreateProject(t *testing.T) {
	db := mocks.NewService(t)
	s := router.NewRouter(db)
	s.GenerateRoutes()
	testURL := "/api/v2/projects"
	testData := []struct {
		name          string
		form          map[string]any
		status        int
		body          map[string]any
		createProject *mocks.Service_CreateProject_Call
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
			name:   "should pass a name",
			form:   map[string]any{},
			status: http.StatusBadRequest,
			body: map[string]any{
				"name":       "El nombre es obligatorio",
				"is_active":  "El estado del projecto es obligatorio",
				"gross_area": "El área bruta es obligatorio",
				"net_area":   "El área neta es obligatorio",
			},
		},
		{
			name: "should pass a status",
			form: map[string]any{
				"name": "test",
			},
			status: http.StatusBadRequest,
			body: map[string]any{
				"is_active":  "El estado del projecto es obligatorio",
				"gross_area": "El área bruta es obligatorio",
				"net_area":   "El área neta es obligatorio",
			},
		},
		{
			name: "should pass a gross area",
			form: map[string]any{
				"name":      "test",
				"is_active": true,
			},
			status: http.StatusBadRequest,
			body: map[string]any{
				"gross_area": "El área bruta es obligatorio",
				"net_area":   "El área neta es obligatorio",
			},
		},
		{
			name: "should pass a net area",
			form: map[string]any{
				"name":       "test",
				"is_active":  true,
				"gross_area": 10,
			},
			status: http.StatusBadRequest,
			body: map[string]any{
				"net_area": "El área neta es obligatorio",
			},
		}, {
			name: "should crate a project",
			form: map[string]any{
				"name":       "test",
				"is_active":  true,
				"gross_area": 10,
				"net_area":   10,
			},
			status: http.StatusCreated,
			body: map[string]any{
				"message": "Proyecto creado correctamente",
			},
			createProject: db.EXPECT().CreateProject(types.Project{
				Name:      "test",
				IsActive:  true,
				GrossArea: 10,
				NetArea:   10,
			}).Return(nil),
		},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			var read io.Reader = nil

			if tt.createProject != nil {
				tt.createProject.Times(1)
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

			for k, v := range tt.body {
				assert.Equal(t, v, mapBody[k])
			}
		})
	}
}
