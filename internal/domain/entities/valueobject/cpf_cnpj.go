package valueobject

import "regexp"

type CPF_CNPJ string

var (
	cpfRegex  = regexp.MustCompile(`^\d{3}\.?\d{3}\.?\d{3}-?\d{2}$`)
	cnpjRegex = regexp.MustCompile(`^\d{2}\.?\d{3}\.?\d{3}/?\d{4}-?\d{2}$`)
)

//crie um parse CPF_CNPJ que retorne um CPF_CNPJ
func ParseCPF_CNPJ(value string) CPF_CNPJ {
	return CPF_CNPJ(value)
}

func (v CPF_CNPJ) IsValidFormat() bool {
	return cpfRegex.MatchString(string(v)) || cnpjRegex.MatchString(string(v))
}
func (v CPF_CNPJ) IsCPF() bool {
	return cpfRegex.MatchString(string(v))
}

func (v CPF_CNPJ) IsCNPJ() bool {
	return cnpjRegex.MatchString(string(v))
}

func (v CPF_CNPJ) String() string {
	return string(v)
}
