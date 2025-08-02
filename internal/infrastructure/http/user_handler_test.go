package http_test

import (
	"encoding/json"
	"mecanica_xpto/internal/domain/user-example"
	"mecanica_xpto/internal/domain/user-example/repository"
	"mecanica_xpto/internal/http"
	http2 "mecanica_xpto/internal/infrastructure/http"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestUserHandler_GetUser(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	repo := repository.NewMemoryRepository()

	// Add a test user-example to the repository
	testUser := &user_example.User{
		ID:    "1",
		Name:  "Test User",
		Email: "test@example.com",
	}
	repo.AddUser(testUser)

	h := http2.NewUserHandler(repo)

	// Create a new Gin engine for testing
	r := gin.New()
	r.GET("/users/:id", h.GetUser)

	tests := []struct {
		name           string
		userID         string
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:           "Existing user-example",
			userID:         "1",
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"id":    "1",
				"name":  "Test User",
				"email": "test@example.com",
			},
		},
		{
			name:           "Non-existing user-example",
			userID:         "999",
			expectedStatus: http.StatusNotFound,
			expectedBody: map[string]interface{}{
				"error": "user-example not found",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/users/"+tt.userID, nil)
			r.ServeHTTP(w, req)

			// Assert status code
			assert.Equal(t, tt.expectedStatus, w.Code)

			// Parse response body
			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			// Assert response body
			assert.Equal(t, tt.expectedBody, response)
		})
	}
}
