package valueobject

type Password string

func ParsePassword(value string) Password {
	return Password(value)
}

func (p Password) String() string {
	return string(p)
}

func (p Password) IsValid() bool {
	// A simple password validation: at least 8 characters, one uppercase, one lowercase, one digit
	if len(p) < 8 {
		return false
	}

	hasUpper := false
	hasLower := false
	hasDigit := false

	for _, char := range p {
		switch {
		case char >= 'A' && char <= 'Z':
			hasUpper = true
		case char >= 'a' && char <= 'z':
			hasLower = true
		case char >= '0' && char <= '9':
			hasDigit = true
		}
	}

	return hasUpper && hasLower && hasDigit
}

func (p Password) IsEmpty() bool {
	return len(p) == 0
}

func (p Password) IsEqual(other Password) bool {
	return p.String() == other.String()
}

func (p Password) IsNotEqual(other Password) bool {
	return !p.IsEqual(other)
}

func (p Password) IsSame(other Password) bool {
	return p.IsEqual(other)
}

// IsValidFormat checks if the password meets the criteria for a valid format.
func (p Password) IsValidFormat() bool {
	// This method can be used to check if the password meets specific format requirements.
	// For now, we can assume that if it is valid, it has already been checked by IsValid.
	return p.IsValid()
}
