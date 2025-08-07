package http_test

//
//import (
//	"bytes"
//	"encoding/json"
//	"mecanica_xpto/internal/domain/repository/user-example"
//	"mecanica_xpto/internal/domain/repository/user-example/repository"
//	http2 "mecanica_xpto/internal/infrastructure/http"
//	"net/http"
//	"net/http/httptest"
//	"testing"
//
//	"github.com/gin-gonic/gin"
//	"github.com/stretchr/testify/assert"
//)
//
//func TestGetUser(t *testing.T) {
//	gin.SetMode(gin.TestMode)
//
//	t.Run("success", func(t *testing.T) {
//		repo := repository.NewMemoryRepository()
//		handler := http2.NewUserHandler(repo)
//
//		testUser := &user_example.User{
//			ID:    "1",
//			Name:  "Test User",
//			Email: "test@example.com",
//		}
//		repo.AddUser(testUser)
//
//		w := httptest.NewRecorder()
//		c, _ := gin.CreateTestContext(w)
//		c.Params = gin.Params{{Key: "id", Value: "1"}}
//
//		handler.GetUser(c)
//
//		assert.Equal(t, http.StatusOK, w.Code)
//		var response map[string]interface{}
//		err := json.Unmarshal(w.Body.Bytes(), &response)
//		assert.NoError(t, err)
//		assert.Equal(t, testUser.ID, response["id"])
//		assert.Equal(t, testUser.Name, response["name"])
//		assert.Equal(t, testUser.Email, response["email"])
//	})
//
//	t.Run("user not found", func(t *testing.T) {
//		repo := repository.NewMemoryRepository()
//		handler := http2.NewUserHandler(repo)
//
//		w := httptest.NewRecorder()
//		c, _ := gin.CreateTestContext(w)
//		c.Params = gin.Params{{Key: "id", Value: "999"}}
//
//		handler.GetUser(c)
//
//		assert.Equal(t, http.StatusNotFound, w.Code)
//	})
//
//	t.Run("invalid id format", func(t *testing.T) {
//		repo := repository.NewMemoryRepository()
//		handler := http2.NewUserHandler(repo)
//
//		w := httptest.NewRecorder()
//		c, _ := gin.CreateTestContext(w)
//		c.Params = gin.Params{{Key: "id", Value: ""}}
//
//		handler.GetUser(c)
//
//		assert.Equal(t, http.StatusBadRequest, w.Code)
//	})
//}
//
//func TestCreateUser(t *testing.T) {
//	gin.SetMode(gin.TestMode)
//
//	t.Run("success", func(t *testing.T) {
//		repo := repository.NewMemoryRepository()
//		handler := http2.NewUserHandler(repo)
//
//		newUser := &user_example.User{
//			Name:  "New User",
//			Email: "new@example.com",
//		}
//
//		userJSON, _ := json.Marshal(newUser)
//		w := httptest.NewRecorder()
//		c, _ := gin.CreateTestContext(w)
//		c.Request = httptest.NewRequest("POST", "/", bytes.NewBuffer(userJSON))
//		c.Request.Header.Set("Content-Type", "application/json")
//
//		handler.CreateUser(c)
//
//		assert.Equal(t, http.StatusCreated, w.Code)
//		var response map[string]interface{}
//		err := json.Unmarshal(w.Body.Bytes(), &response)
//		assert.NoError(t, err)
//		assert.NotEmpty(t, response["id"])
//		assert.Equal(t, newUser.Name, response["name"])
//		assert.Equal(t, newUser.Email, response["email"])
//	})
//
//	t.Run("invalid request body", func(t *testing.T) {
//		repo := repository.NewMemoryRepository()
//		handler := http2.NewUserHandler(repo)
//
//		w := httptest.NewRecorder()
//		c, _ := gin.CreateTestContext(w)
//		c.Request = httptest.NewRequest("POST", "/", bytes.NewBuffer([]byte("invalid json")))
//		c.Request.Header.Set("Content-Type", "application/json")
//
//		handler.CreateUser(c)
//
//		assert.Equal(t, http.StatusBadRequest, w.Code)
//	})
//
//	t.Run("missing required fields", func(t *testing.T) {
//		repo := repository.NewMemoryRepository()
//		handler := http2.NewUserHandler(repo)
//
//		invalidUser := &user_example.User{
//			// Missing required fields
//		}
//
//		userJSON, _ := json.Marshal(invalidUser)
//		w := httptest.NewRecorder()
//		c, _ := gin.CreateTestContext(w)
//		c.Request = httptest.NewRequest("POST", "/", bytes.NewBuffer(userJSON))
//		c.Request.Header.Set("Content-Type", "application/json")
//
//		handler.CreateUser(c)
//
//		assert.Equal(t, http.StatusBadRequest, w.Code)
//	})
//}
//
//func TestUpdateUser(t *testing.T) {
//	gin.SetMode(gin.TestMode)
//
//	t.Run("success", func(t *testing.T) {
//		repo := repository.NewMemoryRepository()
//		handler := http2.NewUserHandler(repo)
//
//		// First create a user
//		existingUser := &user_example.User{
//			ID:    "1",
//			Name:  "Original Name",
//			Email: "original@example.com",
//		}
//		repo.AddUser(existingUser)
//
//		// Update the user
//		updatedUser := &user_example.User{
//			ID:    "1",
//			Name:  "Updated Name",
//			Email: "updated@example.com",
//		}
//
//		userJSON, _ := json.Marshal(updatedUser)
//		w := httptest.NewRecorder()
//		c, _ := gin.CreateTestContext(w)
//		c.Params = gin.Params{{Key: "id", Value: "1"}}
//		c.Request = httptest.NewRequest("PUT", "/", bytes.NewBuffer(userJSON))
//		c.Request.Header.Set("Content-Type", "application/json")
//
//		handler.UpdateUser(c)
//
//		assert.Equal(t, http.StatusOK, w.Code)
//		var response map[string]interface{}
//		err := json.Unmarshal(w.Body.Bytes(), &response)
//		assert.NoError(t, err)
//		assert.Equal(t, updatedUser.Name, response["name"])
//		assert.Equal(t, updatedUser.Email, response["email"])
//	})
//
//	t.Run("user not found", func(t *testing.T) {
//		repo := repository.NewMemoryRepository()
//		handler := http2.NewUserHandler(repo)
//
//		updatedUser := &user_example.User{
//			ID:    "999",
//			Name:  "Updated Name",
//			Email: "updated@example.com",
//		}
//
//		userJSON, _ := json.Marshal(updatedUser)
//		w := httptest.NewRecorder()
//		c, _ := gin.CreateTestContext(w)
//		c.Params = gin.Params{{Key: "id", Value: "999"}}
//		c.Request = httptest.NewRequest("PUT", "/", bytes.NewBuffer(userJSON))
//		c.Request.Header.Set("Content-Type", "application/json")
//
//		handler.UpdateUser(c)
//
//		assert.Equal(t, http.StatusNotFound, w.Code)
//	})
//
//	t.Run("invalid request body", func(t *testing.T) {
//		repo := repository.NewMemoryRepository()
//		handler := http2.NewUserHandler(repo)
//
//		w := httptest.NewRecorder()
//		c, _ := gin.CreateTestContext(w)
//		c.Params = gin.Params{{Key: "id", Value: "1"}}
//		c.Request = httptest.NewRequest("PUT", "/", bytes.NewBuffer([]byte("invalid json")))
//		c.Request.Header.Set("Content-Type", "application/json")
//
//		handler.UpdateUser(c)
//
//		assert.Equal(t, http.StatusBadRequest, w.Code)
//	})
//}
//
//func TestDeleteUser(t *testing.T) {
//	gin.SetMode(gin.TestMode)
//
//	t.Run("success", func(t *testing.T) {
//		repo := repository.NewMemoryRepository()
//		handler := http2.NewUserHandler(repo)
//
//		// First create a user
//		existingUser := &user_example.User{
//			ID:    "1",
//			Name:  "Test User",
//			Email: "test@example.com",
//		}
//		repo.AddUser(existingUser)
//
//		w := httptest.NewRecorder()
//		c, _ := gin.CreateTestContext(w)
//		c.Params = gin.Params{{Key: "id", Value: "1"}}
//
//		handler.DeleteUser(c)
//
//		assert.Equal(t, http.StatusNoContent, w.Code)
//	})
//
//	t.Run("user not found", func(t *testing.T) {
//		repo := repository.NewMemoryRepository()
//		handler := http2.NewUserHandler(repo)
//
//		w := httptest.NewRecorder()
//		c, _ := gin.CreateTestContext(w)
//		c.Params = gin.Params{{Key: "id", Value: "999"}}
//
//		handler.DeleteUser(c)
//
//		assert.Equal(t, http.StatusNotFound, w.Code)
//	})
//
//	t.Run("invalid id format", func(t *testing.T) {
//		repo := repository.NewMemoryRepository()
//		handler := http2.NewUserHandler(repo)
//
//		w := httptest.NewRecorder()
//		c, _ := gin.CreateTestContext(w)
//		c.Params = gin.Params{{Key: "id", Value: ""}}
//
//		handler.DeleteUser(c)
//
//		assert.Equal(t, http.StatusBadRequest, w.Code)
//	})
//}
