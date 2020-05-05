package main

import (
	"log"

	"github.com/nasa9084/gac/commands"
)

func main() {
	if err := commands.Run(); err != nil {
		log.Fatal(err)
	}
}
