package usecase

import (
	"mecanica_xpto/internal/domain/model/dto"
	"mecanica_xpto/internal/domain/model/valueobject"
	"mecanica_xpto/internal/domain/repository/users"
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
	Login(authDTO dto.AuthDTO) (string, *pkg.AppError)
}

type authUseCase struct {
	jwtService *utils.JWTService
	userRepo   users.IUserRepository
}

func NewAuthUseCase(jwtService *utils.JWTService, userRepo users.IUserRepository) *authUseCase {
	return &authUseCase{
		jwtService: jwtService,
		userRepo:   userRepo,
	}
}

// Login handles user login and returns a JWT token
func (a *authUseCase) Login(authDTO dto.AuthDTO) (string, *pkg.AppError) {
	userFromDB, err := a.userRepo.GetByEmail(authDTO.Email)
	if err != nil {
		return "", pkg.NewDomainErrorSimple(ErrCodeInvalidCredential, ErrMsgInvalidCredential, http.StatusUnauthorized)
	}

	hashedPass := valueobject.Password(userFromDB.Password)

	if !hashedPass.Verify(authDTO.Password) {
		return "", pkg.NewDomainErrorSimple(ErrCodeInvalidCredential, ErrMsgInvalidCredential, http.StatusUnauthorized)
	}

	token, err := a.jwtService.GenerateToken(userFromDB.Email)
	if err != nil {
		return "", pkg.NewInfraError(ErrCodeTokenGeneration, ErrMsgTokenGeneration, err, http.StatusInternalServerError)
	}

	return token, nil
}
