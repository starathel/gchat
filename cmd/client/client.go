package main

import (
	"log"

	"github.com/starathel/gchat/internal/client/ui"
)

func main() {
	if err := ui.StartBubbleTea(); err != nil {
		log.Fatal(err)
	}
}
