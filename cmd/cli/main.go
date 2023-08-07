package main

import (
	"encoding/base64"
	"fmt"
	"gits/internal/utils"
)

const password = "Admin1234"

func main() {
	hashPassword, err := utils.HashPassword(password)
	if err != nil {
		_ = fmt.Errorf("can't hash password")
		return
	}
	hashPasswordBytes := []byte(*hashPassword)
	base64HashPass := base64.StdEncoding.EncodeToString(hashPasswordBytes)
	fmt.Printf("password for admin: %v", base64HashPass)
}
