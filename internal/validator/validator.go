package validator

import (
	"regexp"
)

type Validator struct {
	Errors map[string]string
}

func ValidateMail(email string) bool {
	// Expresión regular para validar el email
	regex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(regex)
	return re.MatchString(email)
}
func New() *Validator {
	return &Validator{Errors: make(map[string]string)}
}

func (v *Validator) Valid() bool {
	return len(v.Errors) < 1
}

func (v *Validator) AddError(key, msg string) {
	v.Errors[key] = msg
}

func (v *Validator) Check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

func PermittedValue[T comparable](input T, permitted []T) bool {
	for _, v := range permitted {
		if v == input {
			return true
		}
	}
	return false
}
