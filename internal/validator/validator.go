package validator

import (
	"slices"
	"strings"
	"unicode/utf8"
)

// contains map of *validation error messages* for our form fields
type Validator struct {
	FieldErrors map[string]string
}

func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0
}

func (v *Validator) AddFieldError(key, msg string) {
	// initialize the error map if it's not already.
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}

	if _, exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = msg
	}
}

// if validation check fails, add the error message to the corresponding key (form field).
func (v *Validator) CheckField(check bool, key, msg string) {
	if !check {
		v.AddFieldError(key, msg)
	}
}

func NotBlank(val string) bool {
	return len(strings.TrimSpace(val)) > 0
}

func MaxChars(val string, n int) bool {
	return utf8.RuneCountInString(val) <= n
}

func MinChars(val string, n int) bool {
	return utf8.RuneCountInString(val) >= n
}

// returns true if `val` contains n bytes or less
func MaxBytes(val string, n int) bool {
	return len(val) <= n
}

func PermittedValue[T comparable](val T, permittedValues ...T) bool {
	return slices.Contains(permittedValues, val)
}
