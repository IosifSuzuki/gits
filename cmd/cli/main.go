package main

import (
	"encoding/base64"
	"fmt"
	"github.com/spf13/cobra"
	"gits/internal/utils"
)

var generatePasswordCmd = cobra.Command{
	Use:   "pass",
	Short: "generate hash password by work",
	Long:  ``,
	Run:   generatePassword,
}

func main() {
	if err := generatePasswordCmd.Execute(); err != nil {
		fmt.Printf("error: command pass has failed %v", err)
		return
	}
}

func generatePassword(_ *cobra.Command, args []string) {
	if len(args) <= 1 {
		fmt.Printf("error: please specify a password")
		return
	}
	password := args[1]
	hashPassword, err := utils.HashPassword(password)
	if err != nil {
		fmt.Printf("error: generate hash password has failed with error: %v", err)
		return
	} else if hashPassword == nil {
		fmt.Printf("error: ash password has nil value")
		return
	}
	hashPasswordBytes := []byte(*hashPassword)
	base64HashPass := base64.StdEncoding.EncodeToString(hashPasswordBytes)
	fmt.Printf("hash password has generated: %v", base64HashPass)
}
