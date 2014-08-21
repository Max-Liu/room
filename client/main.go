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

	client, err := link.Dial("tcp", "127.0.0.1:10010", protocol)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	client.OnMessage(func(session *link.Session, message []byte) {
		println("message:", string(message))
	})

	client.OnClose(func(session *link.Session, reason error) {
		println("closed")
	})

	client.Start()

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

	stream := room.NewStream(room.Box{user, "user"})
	client.Send(stream)

	go func() {
		for {
			<-time.Tick(1 * time.Second)
			stream = room.NewStream(room.Box{"Sending Ticker~~~", "debug"})
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
		stream := room.NewStream(room.Box{user, "user"})

		client.Send(stream)
	}
	client.Close(nil)

	println("bye")
}
