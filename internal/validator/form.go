package validator

import (
	"github.com/go-playground/validator/v10"
	"regexp"
	"strconv"
	"strings"
)

const (
	MinCountUsername = 4
	MinCountPassword = 8
	MinCountFullName = 6
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

func FullNameValidator(fieldLevel validator.FieldLevel) bool {
	value, ok := fieldLevel.Field().Interface().(string)
	if !ok {
		return false
	}
	if len(strings.Split(value, " ")) < 2 {
		return false
	}
	if len(value) < MinCountFullName {
		return false
	}

	return true
}

func FilledIn(fieldLevel validator.FieldLevel) bool {
	value, ok := fieldLevel.Field().Interface().(string)
	if !ok {
		return false
	}

	minValue, _ := strconv.Atoi(fieldLevel.Param())
	value = strings.TrimSpace(value)

	return len(value) >= minValue
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
