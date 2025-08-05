package middleware

import (
	"mecanica_xpto/internal/infrastructure/config"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestJWTAuthMiddleware(t *testing.T) {
	// Carregar a chave secreta do .env
	secretKey := config.GetSecretKey()

	gin.SetMode(gin.TestMode)
	r := gin.New()

	// Simulando um handler para testar o middleware
	r.GET("/test", JWTAuthMiddleware(secretKey), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Access Granted"})
	})

	// Teste 1: Requisição sem token
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	// Verificando se a resposta é a esperada quando o token não é fornecido
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Missing authorization token!")

	// Teste 2: Requisição com token inválido (formato errado)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "InvalidToken")
	r.ServeHTTP(w, req)

	// Verificando se a resposta é a esperada quando o token tem formato errado
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid token format!")

	// Teste 3: Requisição com token válido
	validToken, _ := GenerateToken("testuser", secretKey)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+validToken)
	r.ServeHTTP(w, req)

	// Verificando se a resposta é a esperada quando o token é válido
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Access Granted")

	// Teste 4: Requisição com token inválido (chave secreta errada)
	wrongSecretKey := "wrongsecret"
	invalidToken, _ := GenerateToken("testuser", wrongSecretKey)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+invalidToken)
	r.ServeHTTP(w, req)

	// Verificando se a resposta é a esperada quando a chave secreta está errada
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "signature is invalid")
}

func TestGenerateToken(t *testing.T) {
	// Testando se a geração de token está funcionando corretamente
	secretKey := config.GetSecretKey() // Pega a chave secreta do .env
	userID := "testuser"

	tokenString, err := GenerateToken(userID, secretKey)
	assert.NoError(t, err)
	assert.NotEmpty(t, tokenString)

	// Validar o token gerado
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	assert.NoError(t, err)
	assert.True(t, token.Valid)

	// Verificando os dados do token
	claims, ok := token.Claims.(*Claims)
	assert.True(t, ok)
	assert.Equal(t, userID, claims.UserID)
}
