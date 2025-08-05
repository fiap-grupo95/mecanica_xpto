package valueobject

import "regexp"

type CpfCnpj string

var (
	cpfRegex  = regexp.MustCompile(`^\d{3}\.?\d{3}\.?\d{3}-?\d{2}$`)
	cnpjRegex = regexp.MustCompile(`^\d{2}\.?\d{3}\.?\d{3}/?\d{4}-?\d{2}$`)
)

// ParseCPF_CNPJ crie um parse CpfCnpj que retorne um CpfCnpj
func ParseCPF_CNPJ(value string) CpfCnpj {
	return CpfCnpj(value)
}

func (v CpfCnpj) IsValidFormat() bool {
	return cpfRegex.MatchString(string(v)) || cnpjRegex.MatchString(string(v))
}
func (v CpfCnpj) IsCPF() bool {
	return cpfRegex.MatchString(string(v))
}

func (v CpfCnpj) IsCNPJ() bool {
	return cnpjRegex.MatchString(string(v))
}

func (v CpfCnpj) String() string {
	return string(v)
}
