package utils

import (
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
)

func HashPassword(password string) (*string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return nil, err
	}
	hashStr := string(hash)
	return &hashStr, nil
}

func CompareHashAndPassword(hash string, password string) (bool, error) {
	lhsBytes := []byte(hash)
	rhsBytes := []byte(password)
	err := bcrypt.CompareHashAndPassword(lhsBytes, rhsBytes)
	if err != nil {
		return false, err
	}
	return true, err
}

func DateFormat(time time.Time, format string) string {
	return time.Format(format)
}

func Add(lsh, rhs int) int {
	return lsh + rhs
}

func PrefixString(text string, words int) string {
	parts := strings.Split(text, " ")
	var actualWords int
	if len(parts)-1 > words {
		actualWords = words
	} else {
		actualWords = len(parts) - 1
	}
	return strings.Join(parts[:actualWords], " ")
}
