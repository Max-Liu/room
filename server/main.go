package main

import (
	"log"
	"room"
)

func main() {
	log.SetFlags(log.Llongfile)
	rpcListen := room.NewRpcListener()

	go func() {
		room.RpcWorker(rpcListen)
	}()

	go func() {
		debugRoom := room.NewDebugChatRoom()
		debugRoom.Start()
		rpcListen.Msg <- <-debugRoom.Msg

	}()

	for {
		<-rpcListen.CreateNewRoom
		go func() {
			newRoom := room.NewChatRoom()
			newRoom.Start()
			rpcListen.Msg <- <-newRoom.Msg
		}()
	}

}
