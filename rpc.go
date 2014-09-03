package room

import (
	"log"
	"strconv"
	"strings"
)

type RpcListener struct {
	CreateNewRoom  chan int
	HasCreatedRoom chan int
	Msg            []byte
}

func NewRpcListener() *RpcListener {
	createNewRoom := make(chan int)
	hasCreatedRoom := make(chan int)

	var Msg []byte

	return &RpcListener{
		createNewRoom,
		hasCreatedRoom,
		Msg,
	}

}

func (l *RpcListener) CreateChatRoom(line []byte, reply *[]byte) error {
	l.CreateNewRoom <- 1
	<-l.HasCreatedRoom
	*reply = l.Msg
	return nil
}
func (l *RpcListener) EndChatRoom(line []byte, reply *[]byte) error {
	endPortInt, err := strconv.Atoi(string(line))

	if err != nil {
		log.Fatal(err)
	}

	for _, server := range RegRoomList {
		serverPortStr := strings.Split(server.Listener().Addr().String(), ":")[1]
		serverPortInt, err := strconv.Atoi(serverPortStr)
		if err != nil {
			log.Fatal(err)
		}
		if serverPortInt == endPortInt {
			server.Stop("end")
		}
	}

	return nil
}
