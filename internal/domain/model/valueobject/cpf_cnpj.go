package valueobject

import (
	"mecanica_xpto/pkg/utils"
	"mecanica_xpto/pkg/validators"
	"regexp"
)

type CpfCnpj string

const (
	sizeCNPJWithoutDV       = 12
	baseValue               = int('0')
	cnpjZerosOnly           = "00000000000000"
	operationAndErrorFormat = "%s: %w"
)

var (
	regexCPFPattern      = regexp.MustCompile("^[0-9]{11}$")
	regexCNPJWithoutDV   = regexp.MustCompile("^[0-9]{12}$")
	regexCNPJ            = regexp.MustCompile("^[0-9]{14}$")
	regexCharsNotAllowed = regexp.MustCompile("[^0-9./-]")
	regexNonNumericChars = regexp.MustCompile("[^0-9]")
	weightDV             = []int{6, 5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}
)

func NewCpfCnpj(v string) (CpfCnpj, error) {
	v = clean(v)
	c := CpfCnpj(v)
	if err := c.IsValid(); err != nil {
		return "", err
	}
	return c, nil
}

func (c CpfCnpj) IsValid() error {
	if len(c) <= 11 {
		return validators.CnpjIsValid(c.String())
	}
	return validators.CnpjIsValid(c.String())
}

func (c CpfCnpj) Mask() string {
	if len(c) == 11 {
		return utils.MaskCPF(c.String())
	}
	return utils.MaskCNPJ(c.String())
}

func (c CpfCnpj) String() string {
	return string(c)
}

func clean(c string) string {
	return regexNonNumericChars.ReplaceAllString(c, "")
}
