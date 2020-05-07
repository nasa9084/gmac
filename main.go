package main

import (
	"github.com/nasa9084/gmac/commands"
	"github.com/nasa9084/gmac/log"
)

func main() {
	if err := commands.Run(); err != nil {
		log.Fatal(err)
	}
}
