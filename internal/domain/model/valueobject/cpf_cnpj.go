package valueobject

import "regexp"

type Document string

var (
	cpfRegex  = regexp.MustCompile(`^\d{3}\.?\d{3}\.?\d{3}-?\d{2}$`)
	cnpjRegex = regexp.MustCompile(`^\d{2}\.?\d{3}\.?\d{3}/?\d{4}-?\d{2}$`)
)

// ParseCPF_CNPJ crie um parse CpfCnpj que retorne um CpfCnpj
func ParseDocument(value string) Document {
	return Document(value)
}

func (v Document) IsValidFormat() bool {
	return cpfRegex.MatchString(string(v)) || cnpjRegex.MatchString(string(v))
}
func (v Document) IsCPF() bool {
	return cpfRegex.MatchString(string(v))
}

func (v Document) IsCNPJ() bool {
	return cnpjRegex.MatchString(string(v))
}

func (v Document) String() string {
	return string(v)
}
