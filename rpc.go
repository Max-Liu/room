package room

import (
	"encoding/json"
	"log"
	"net"
	"net/rpc"
	"strconv"
)

type RpcListener struct {
	CreateNewRoom chan int
	Msg           chan []byte
}

func NewRpcListener() *RpcListener {
	createNewRoom := make(chan int)
	Msg := make(chan []byte)

	return &RpcListener{
		createNewRoom,
		Msg,
	}

}
func RpcWorker(l *RpcListener) {
	addy, err := net.ResolveTCPAddr("tcp", "127.0.0.1:42586")
	if err != nil {
		log.Fatal(err)
	}

	inbound, err := net.ListenTCP("tcp", addy)
	if err != nil {
		log.Fatal(err)
	}
	rpc.Register(l)
	rpc.Accept(inbound)
}

func (l *RpcListener) GetRoomNum(line []byte, reply *[]byte) error {
	RoomNumStr := strconv.Itoa(len(RegRoomList))
	*reply = []byte(RoomNumStr)
	return nil
}

func (l *RpcListener) Ping(line []byte, reply *[]byte) error {
	*reply = []byte("pong~~~")
	return nil
}

func (l *RpcListener) GetRoomList(line []byte, reply *[]byte) error {
	var roomInfoList []RoomSysInfo
	for _, RegRoom := range RegRoomList {
		var roomInfo RoomSysInfo
		roomInfo.Addr = RegRoom.State.(RoomSysInfo).Addr
		roomInfo.Pid = RegRoom.State.(RoomSysInfo).Pid
		roomInfo.StartTime = RegRoom.State.(RoomSysInfo).StartTime
		roomInfoList = append(roomInfoList, roomInfo)
	}
	jsonByte, err := json.Marshal(roomInfoList)
	*reply = jsonByte
	return err
}

func (l *RpcListener) CreateChatRoom(line []byte, reply *[]byte) error {
	l.CreateNewRoom <- 1
	*reply = <-l.Msg
	return nil
}
func (l *RpcListener) EndChatRoom(line []byte, reply *[]byte) error {
	endPortInt, err := strconv.Atoi(string(line))

	if err != nil {
		log.Fatal(err)
	}

	RegRoomList[endPortInt].Stop("end")
	delete(RegRoomList, endPortInt)

	return nil
}
