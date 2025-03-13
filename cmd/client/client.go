package main

import (
	"fmt"
	"log"

	"github.com/starathel/gchat/internal/client/transport"
)

func main() {
	//if err := ui.StartBubbleTea(); err != nil {
	//	log.Fatal(err)
	//}

	client, err := transport.NewClient()
	if err != nil {
		log.Fatal(err)
	}

	rooms, err := client.GetRoomsList()
	if err != nil {
		log.Fatal(err)
	}

	for _, room := range rooms {
		fmt.Printf("%s %d\n", room.Id, room.UsersCount)
	}
}
