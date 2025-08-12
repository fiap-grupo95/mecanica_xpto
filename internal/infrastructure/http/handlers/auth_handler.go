package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"mecanica_xpto/internal/domain/model/dto"
	"mecanica_xpto/internal/domain/usecase"
	"mecanica_xpto/pkg"
)

const (
	ErrCodeInvalidRequest = "INVALID_REQUEST"
	ErrMsgInvalidRequest  = "Invalid request body"
)

type AuthHandler struct {
	usecase usecase.AuthInterface
}

func NewAuthHandler(usecase usecase.AuthInterface) *AuthHandler {
	return &AuthHandler{
		usecase: usecase,
	}
}

// Login autentica um usuário e retorna um token JWT.
//
// @Summary      Autenticação do usuário
// @Description  Autentica um usuário com email e senha e retorna token JWT.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        loginRequest  body  dto.AuthDTO  true  "Credenciais do usuário"
// @Success      200  {object}  map[string]string  "Token JWT"
// @Failure      400  {object}  pkg.ErrorResponse  "Requisição inválida"
// @Failure      401  {object}  pkg.ErrorResponse  "Não autorizado"
// @Router       /v1/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.AuthDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(
			http.StatusBadRequest,
			pkg.NewDomainErrorSimple(ErrCodeInvalidRequest, ErrMsgInvalidRequest, http.StatusBadRequest).ToHTTPError(),
		)
		return
	}

	token, errLogin := h.usecase.Login(req)
	if errLogin != nil {
		c.JSON(http.StatusUnauthorized, errLogin.ToHTTPError())
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
