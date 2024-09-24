package validator

type Validator struct {
	Errors map[string]string
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
