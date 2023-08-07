package utils

import "golang.org/x/crypto/bcrypt"

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
