package room

import (
	"encoding/json"
	"log"
	"strconv"
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

func (l *RpcListener) GetRoomNum(line []byte, reply *[]byte) error {
	RoomNumStr := strconv.Itoa(len(RegRoomList))
	*reply = []byte(RoomNumStr)
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
	<-l.HasCreatedRoom
	*reply = l.Msg
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
