package handler_test

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"mecanica_xpto/internal/domain/user"
	"mecanica_xpto/internal/domain/user/repository"
	"mecanica_xpto/internal/handler"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserHandler_GetUser(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	repo := repository.NewMemoryRepository()

	// Add a test user to the repository
	testUser := &user.User{
		ID:    "1",
		Name:  "Test User",
		Email: "test@example.com",
	}
	repo.AddUser(testUser)

	h := handler.NewUserHandler(repo)

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
			name:           "Existing user",
			userID:         "1",
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"id":    "1",
				"name":  "Test User",
				"email": "test@example.com",
			},
		},
		{
			name:           "Non-existing user",
			userID:         "999",
			expectedStatus: http.StatusNotFound,
			expectedBody: map[string]interface{}{
				"error": "user not found",
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
