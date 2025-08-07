package valueobject

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
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
		return cpfIsValid(c.String())
	}
	return cnpjIsValid(c.String())
}

func (c CpfCnpj) Mask() string {
	if len(c) == 11 {
		return maskCPF(c.String())
	}
	return maskCNPJ(c.String())
}

func (c CpfCnpj) String() string {
	return string(c)
}

func cpfIsValid(c string) error {
	if !regexCPFPattern.MatchString(c) {
		return errors.New("invalid CPF format")
	}

	charsRepeated := true
	for i := 1; i < len(c); i++ {
		if c[0] != c[i] {
			charsRepeated = false
			break
		}
	}
	if charsRepeated {
		return errors.New("CPF cannot have all characters the same")
	}

	i, sum := 10, 0
	for index := 0; index < len(c)-2; index++ {
		pos, _ := strconv.Atoi(string(c[index]))
		sum += pos * i
		i--
	}

	prod := sum * 10
	mod := prod % 11

	if mod == 10 {
		mod = 0
	}

	digit1, _ := strconv.Atoi(string(c[9]))
	if mod != digit1 {
		return errors.New("CPF validation digit 1 does not match")
	}

	i, sum = 11, 0
	for index := 0; index < len(c)-1; index++ {
		pos, _ := strconv.Atoi(string(c[index]))
		sum += pos * i
		i--
	}

	prod = sum * 10
	mod = prod % 11

	if mod == 10 {
		mod = 0
	}

	digit2, _ := strconv.Atoi(string(c[10]))
	if mod != digit2 {
		return errors.New("CPF validation digit 2 does not match")
	}

	return nil
}

func cnpjIsValid(c string) error {
	if regexCharsNotAllowed.MatchString(c) {
		return errors.New("CNPJ contains invalid characters")
	}

	cnpjWithoutMask := clean(c)
	if !regexCNPJ.MatchString(cnpjWithoutMask) || cnpjWithoutMask == cnpjZerosOnly {
		return errors.New("CNPJ has an invalid pattern")
	}

	calculatedCNPJDV, err := calculateAndValidateCNPJDV(cnpjWithoutMask[:sizeCNPJWithoutDV])
	if err != nil {
		return err
	}

	informedCNPJDV := cnpjWithoutMask[sizeCNPJWithoutDV:]
	if informedCNPJDV != calculatedCNPJDV {
		return errors.New("CNPJ validation digits do not match")
	}

	return nil
}

func calculateAndValidateCNPJDV(c string) (string, error) {
	const op = "validator.calculateAndValidateCNPJDV"

	if regexCharsNotAllowed.MatchString(c) {
		return "", fmt.Errorf(
			operationAndErrorFormat,
			op,
			errors.New("CNPJ contains invalid characters"),
		)
	}

	CNPJWithoutMask := clean(c)
	if !regexCNPJWithoutDV.MatchString(CNPJWithoutMask) || CNPJWithoutMask == cnpjZerosOnly[:sizeCNPJWithoutDV] {
		return "", fmt.Errorf(
			operationAndErrorFormat,
			op,
			errors.New("CNPJ has an invalid pattern"),
		)
	}

	sumDV1 := 0
	sumDV2 := 0
	for i := 0; i < sizeCNPJWithoutDV; i++ {
		asciiDigit := int(CNPJWithoutMask[i]) - baseValue

		sumDV1 += asciiDigit * weightDV[i+1]
		sumDV2 += asciiDigit * weightDV[i]
	}

	var dv1 int
	if sumDV1%11 < 2 {
		dv1 = 0
	} else {
		dv1 = 11 - (sumDV1 % 11)
	}

	sumDV2 += dv1 * weightDV[sizeCNPJWithoutDV]

	var dv2 int
	if sumDV2%11 < 2 {
		dv2 = 0
	} else {
		dv2 = 11 - (sumDV2 % 11)
	}

	return fmt.Sprintf("%d%d", dv1, dv2), nil
}

func maskCPF(c string) string {
	re := regexp.MustCompile(`^(\d{3})(\d{3})(\d{3})(\d{2})$`)
	return re.ReplaceAllString(c, "$1.$2.$3-$4")
}

func maskCNPJ(c string) string {
	re := regexp.MustCompile(`^(\d{2})(\d{3})(\d{3})(\d{4})(\d{2})$`)
	return re.ReplaceAllString(c, "$1.$2.$3/$4-$5")
}

func clean(c string) string {
	return regexNonNumericChars.ReplaceAllString(c, "")
}
