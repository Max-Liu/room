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
	msg := make(chan []byte)

	return &ChatRoom{
		SysInfo: RoomSysInfo{},
		Msg:     msg,
	}
}

func (c *ChatRoom) Start() error {
	protocol := link.PacketN(2, binary.BigEndian)
	pid := os.Getpid()

	server, err := link.Listen("tcp", "127.0.0.1:0", protocol)

	c.SysInfo = NewRoomSysInfo(pid, server.Listener().Addr().String(), time.Now().Unix())
	serverPortStr := strings.Split(server.Listener().Addr().String(), ":")[1]
	serverPortInt, _ := strconv.Atoi(serverPortStr)

	RegRoomList[serverPortInt] = server
	box := Box{}
	b, _ := json.Marshal(c.SysInfo)
	fmt.Printf("%s\n", b)
	server.State = c.SysInfo

	go server.AcceptLoop(func(session *link.Session) {
		session.ReadLoop(func(message []byte) {
			json.Unmarshal(message, &box)
			switch box.Kind {
			case "user":
				{
					user := box.Object.(map[string]interface{})
					switch user["CmdContent"] {
					case "reg":
						{
							fmt.Println(user["Name"], "joined the game.")
						}
					case "msg":
						{
							userMsg := user["Msg"].(map[string]interface{})
							fmt.Println(user["Name"], "Say:", userMsg["Content"])
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
	go func() {
		for {
			c.Msg <- b
		}
	}()

	return err
}
