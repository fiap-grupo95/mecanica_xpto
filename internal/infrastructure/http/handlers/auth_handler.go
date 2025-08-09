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

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.UserDTO
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
