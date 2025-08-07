package valueobject

import (
	"encoding/base64"
	"errors"
	"mecanica_xpto/pkg"
	"strings"

	"golang.org/x/crypto/argon2"
)

type Password string // representa sempre uma senha *hashada*

// NewPassword recebe uma senha pura, valida e retorna um Password hashado
func NewPassword(p string) (Password, error) {
	if err := isValid(p); err != nil {
		return "", err
	}

	hashed, err := pkg.HashPassword(p)
	if err != nil {
		return "", err
	}

	return Password(hashed), nil
}

// String retorna o hash como string
func (p Password) String() string {
	return string(p)
}

// Verify compara uma senha pura com o hash armazenado
func (p Password) Verify(plain string) bool {
	parts := strings.Split(p.String(), ".")
	if len(parts) != 2 {
		return false
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[0])
	if err != nil {
		return false
	}

	expectedHash, err := base64.RawStdEncoding.DecodeString(parts[1])
	if err != nil {
		return false
	}

	hashToCompare := argon2.IDKey([]byte(plain), salt, 1, 64*1024, 4, 32)

	return subtleCompare(expectedHash, hashToCompare)
}

// subtleCompare realiza comparação segura
func subtleCompare(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	var result byte
	for i := range a {
		result |= a[i] ^ b[i]
	}
	return result == 0
}

// isValid valida a força mínima da senha pura
func isValid(s string) error {
	if len(s) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	hasUpper := false
	hasLower := false
	hasDigit := false

	for _, char := range s {
		switch {
		case char >= 'A' && char <= 'Z':
			hasUpper = true
		case char >= 'a' && char <= 'z':
			hasLower = true
		case char >= '0' && char <= '9':
			hasDigit = true
		}
	}

	if !(hasUpper && hasLower && hasDigit) {
		return errors.New("password must contain at least one uppercase letter, one lowercase letter, and one digit")
	}
	return nil
}
