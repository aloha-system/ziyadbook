package main

import (
	"log"

	"ziyadbook/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
