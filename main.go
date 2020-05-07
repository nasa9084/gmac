package main

import (
	"log"

	"github.com/nasa9084/gmac/commands"
)

func main() {
	if err := commands.Run(); err != nil {
		log.Fatal(err)
	}
}
