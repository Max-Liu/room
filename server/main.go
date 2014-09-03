package main

import (
	"log"
	"net"
	"net/rpc"
	"room"
)

func main() {
	log.SetFlags(log.Llongfile)
	rpcListen := room.NewRpcListener()

	go func() {
		addy, err := net.ResolveTCPAddr("tcp", "127.0.0.1:42586")
		if err != nil {
			log.Fatal(err)
		}

		inbound, err := net.ListenTCP("tcp", addy)
		if err != nil {
			log.Fatal(err)
		}

		rpc.Register(rpcListen)
		rpc.Accept(inbound)
	}()

	for {
		<-rpcListen.CreateNewRoom
		newRoom := room.NewChatRoom()
		go func() {
			err := newRoom.Start()
			if err != nil {
				log.Fatal(err)
			}
		}()
		rpcListen.Msg = <-newRoom.Msg
		rpcListen.HasCreatedRoom <- 1
	}
}
