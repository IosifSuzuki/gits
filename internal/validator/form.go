package validator

import (
	"github.com/go-playground/validator/v10"
	"regexp"
)

const (
	MinCountUsername = 4
	MinCountPassword = 8
)

func UsernameValidator(fieldLevel validator.FieldLevel) bool {
	patternHasUpperSymbol := regexp.MustCompile("[A-Z]")
	patternHasLowerSymbol := regexp.MustCompile("[a-z]")
	value, ok := fieldLevel.Field().Interface().(string)
	if !ok {
		return false
	}
	if len(value) < MinCountUsername {
		return false
	}
	if !patternHasUpperSymbol.MatchString(value) || !patternHasLowerSymbol.MatchString(value) {
		return false
	}
	return true
}

func PasswordValidator(fieldLevel validator.FieldLevel) bool {
	patternHasUpperSymbol := regexp.MustCompile("[A-Z]")
	patternHasLowerSymbol := regexp.MustCompile("[a-z]")
	patternHasDigitSymbols := regexp.MustCompile(`\d{4,}`)

	value, ok := fieldLevel.Field().Interface().(string)
	if !ok {
		return false
	}
	if len(value) < MinCountPassword {
		return false
	}
	if !patternHasUpperSymbol.MatchString(value) || !patternHasLowerSymbol.MatchString(value) {
		return false
	}
	if !patternHasDigitSymbols.MatchString(value) {
		return false
	}
	return true
}
