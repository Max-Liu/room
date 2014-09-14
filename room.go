package room

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/funny/link"
)

type RoomSysInfo struct {
	Pid       int
	Addr      string
	StartTime int64
}

func NewRoomSysInfo(pid int, addr string, startTime int64) RoomSysInfo {
	return RoomSysInfo{
		pid,
		addr,
		startTime,
	}
}

type ChatRoom struct {
	SysInfo RoomSysInfo
	Msg     chan []byte
}

var RegRoomList map[int]*link.Server

func init() {
	RegRoomList = make(map[int]*link.Server)
}

func NewChatRoom() *ChatRoom {
	msg := make(chan []byte, 1024)

	return &ChatRoom{
		SysInfo: RoomSysInfo{
			0,
			"127.0.0.1:0",
			0,
		},
		Msg: msg,
	}
}

func NewDebugChatRoom() *ChatRoom {

	msg := make(chan []byte, 1024)

	return &ChatRoom{
		SysInfo: RoomSysInfo{
			0,
			"localhost:55000",
			0,
		},
		Msg: msg,
	}

}

func (chatRoom *ChatRoom) Start() {
	protocol := link.PacketN(2, binary.LittleEndian)

	pid := os.Getpid()
	server, err := link.Listen("tcp", chatRoom.SysInfo.Addr, protocol)
	if err != nil {
		log.Println(err)
	}

	chatRoom.SysInfo = NewRoomSysInfo(pid, server.Listener().Addr().String(), time.Now().Unix())
	serverPortStr := strings.Split(server.Listener().Addr().String(), ":")[1]
	serverPortInt, _ := strconv.Atoi(serverPortStr)

	channel := link.NewChannel(server.Protocol())

	RegRoomList[serverPortInt] = server
	b, _ := json.Marshal(chatRoom.SysInfo)
	chatRoom.Msg <- b
	server.State = chatRoom.SysInfo

	go server.AcceptLoop(func(session *link.Session) {
		channel.Join(session, nil)
		session.ReadLoop(func(message []byte) {
			box := Box{}
			json.Unmarshal(message, &box)
			switch box.Kind {
			case "user":
				{
					user := box.Object.(map[string]interface{})
					switch user["CmdContent"] {
					case "reg":
						{
							channel.Broadcast(link.Binary(fmt.Sprintln(user["Name"].(string), "joined the game")))
						}
					case "msg":
						{
							userMsg := user["Msg"].(map[string]interface{})
							channel.Broadcast(link.Binary(fmt.Sprintln(user["Name"], "Say:", userMsg["Content"])))
						}
					}
				}
			case "debug":
				{
					log.Println(box.Object, session.Conn().RemoteAddr().String())
				}
			default:
				{
					log.Println(string(message))
				}
			}
		})
	})
}
