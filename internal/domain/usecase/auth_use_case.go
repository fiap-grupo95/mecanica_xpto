package usecase

import (
	"mecanica_xpto/internal/domain/model/dto"
	"mecanica_xpto/internal/domain/model/entities"
	"mecanica_xpto/internal/domain/model/valueobject"
	"mecanica_xpto/pkg"
	"mecanica_xpto/pkg/utils"
	"net/http"
)

const (
	ErrCodeInvalidCredential = "INVALID_CREDENTIALS"
	ErrMsgInvalidCredential  = "Invalid email or password"
	ErrCodeTokenGeneration   = "TOKEN_GENERATION_ERROR"
	ErrMsgTokenGeneration    = "Failed to generate token"
)

type AuthInterface interface {
	Login(userDTO dto.UserDTO) (string, *pkg.AppError)
}

type authUseCase struct {
	jwtService *utils.JWTService
}

func NewAuthUseCase(jwtService *utils.JWTService) *authUseCase {
	return &authUseCase{
		jwtService: jwtService,
	}
}

// Login handles user login and returns a JWT token
func (a *authUseCase) Login(userDTO dto.UserDTO) (string, *pkg.AppError) {
	// TODO: usar repository de user
	user := entities.User{
		Email:    userDTO.Email,
		Password: userDTO.Password,
	}
	pass, err := valueobject.NewPassword(user.Password)
	if err != nil {
		return "", pkg.NewDomainError(ErrCodeInvalidCredential, ErrMsgInvalidCredential, err, http.StatusBadRequest)
	}
	// TODO: Fim do TODO

	if !pass.Verify(userDTO.Password) {
		return "", pkg.NewDomainErrorSimple(ErrCodeInvalidCredential, ErrMsgInvalidCredential, http.StatusUnauthorized)
	}

	token, err := a.jwtService.GenerateToken(user.Email)
	if err != nil {
		return "", pkg.NewInfraError(ErrCodeTokenGeneration, ErrMsgTokenGeneration, err, http.StatusInternalServerError)
	}

	return token, nil
}
