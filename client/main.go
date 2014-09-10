package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"room"
	"time"

	"github.com/funny/link"
)

func main() {
	log.SetFlags(log.Lshortfile)
	protocol := link.PacketN(2, binary.BigEndian)

	client, err := link.Dial("tcp", "127.0.0.1:52128", protocol)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	go client.ReadLoop(func(message []byte) {
		println("message:", string(message))
	})

	user := room.NewUser()
	fmt.Println("Your Name:")

EnterName:
	if _, err := fmt.Scanf("%s\n", &user.Name); err != nil {
		if user.Name == "" {
			fmt.Println("Please enter your name:")
			goto EnterName
		}
	}

	user.CmdContent = "reg"
	stream := link.JSON{}
	stream.V = room.Box{user, "user"}

	client.Send(stream)

	go func() {
		for {
			<-time.Tick(1 * time.Second)
			stream.V = room.Box{"Sending Ticker~~~", "debug"}
			client.Send(stream)
		}
	}()

	for {
	EnterUserMsg:
		fmt.Println("Say:")
		if _, err := fmt.Scanf("%s\n", &user.Msg.Content); err != nil {
			goto EnterUserMsg
		}
		user.CmdContent = "msg"
		stream.V = room.Box{user, "user"}

		client.Send(stream)
	}
	client.Close(nil)

	println("bye")
}
