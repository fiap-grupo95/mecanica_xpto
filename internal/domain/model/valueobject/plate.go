package valueobject

import "regexp"

type Plate string

var (
	plateRegex = regexp.MustCompile(`^[A-Z]{3}(?:\d[A-Z0-9]\d{2}|\d{4})$`)
)

func ParsePlate(value string) Plate {
	return Plate(value)
}

func (v Plate) IsValidFormat() bool {
	return plateRegex.MatchString(string(v))
}

func (v Plate) String() string {
	return string(v)
}
