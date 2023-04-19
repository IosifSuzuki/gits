package main

import (
	"gits/internal/app"
	"log"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatalln(err)
	}
}
