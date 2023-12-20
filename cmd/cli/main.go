package main

import (
	"encoding/base64"
	"github.com/spf13/cobra"
	"gits/internal/utils"
	"gits/pkg/logger"
	"go.uber.org/zap"
)

var log = logger.NewLogger(logger.DebugLevel)

var generatePasswordCmd = cobra.Command{
	Use:   "pass",
	Short: "generate hash password by work",
	Long:  ``,
	Run:   generatePassword,
}

func main() {
	if err := generatePasswordCmd.Execute(); err != nil {
		log.Error("command failed", zap.Error(err))
		return
	}
}

func generatePassword(_ *cobra.Command, args []string) {
	if len(args) <= 1 {
		log.Error("please specify a password as argument")
		return
	}
	password := args[1]
	hashPassword, err := utils.HashPassword(password)
	if err != nil {
		log.Error("generate hash password has failed", zap.Error(err))
		return
	} else if hashPassword == nil {
		log.Error("hash password has nil value")
		return
	}
	hashPasswordBytes := []byte(*hashPassword)
	base64HashPass := base64.StdEncoding.EncodeToString(hashPasswordBytes)
	log.Info("result", zap.String("base64 hash password", base64HashPass))
}
