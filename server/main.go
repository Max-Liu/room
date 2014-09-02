package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
	"room"
	"strconv"
	"strings"
	"time"

	"github.com/funny/link"
)

var rpcListen *RpcListener
var regRoomList []*link.Server

type RpcListener struct {
	CreateNewRoom  chan int
	HasCreatedRoom chan int
	EndRoom        chan int
	msg            []byte
}

func (l *RpcListener) CreateChatRoom(line []byte, reply *[]byte) error {
	l.CreateNewRoom <- 1
	<-l.HasCreatedRoom
	*reply = l.msg
	return nil
}
func (l *RpcListener) EndChatRoom(line []byte, reply *[]byte) error {
	endPortInt, err := strconv.Atoi(string(line))
	if err != nil {
		log.Fatal(err)
	}
	for _, server := range regRoomList {
		serverPortStr := strings.Split(server.Listener().Addr().String(), ":")[1]
		serverPortInt, err := strconv.Atoi(serverPortStr)
		if err != nil {
			log.Fatal(err)
		}
		if serverPortInt == endPortInt {
			server.Stop()
		}

	}
	return nil
}

func main() {
	log.SetFlags(log.Llongfile)
	rpcListen = new(RpcListener)
	rpcListen.CreateNewRoom = make(chan int)
	rpcListen.HasCreatedRoom = make(chan int)

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

	roomWorker()
}

func roomWorker() {
	for {
		<-rpcListen.CreateNewRoom
		go newRoom()
		rpcListen.HasCreatedRoom <- 1
	}

}

func newRoom() {
	protocol := link.PacketN(2, binary.BigEndian)
	pid := os.Getpid()

	server, err := link.Listen("tcp", "127.0.0.1:0", protocol)

	if err != nil {
		panic(err)
	}

	//println("server start:", server.Listener().Addr().String())

	roomSysInfo := room.NewRoom(pid, server.Listener().Addr().String(), time.Now().Unix())

	regRoomList = append(regRoomList, server)
	box := room.Box{}
	b, _ := json.Marshal(roomSysInfo)
	rpcListen.msg = b
	fmt.Printf("%s\n", b)

	server.Accept(func(session *link.Session) {
		session.OnMessage(func(session *link.Session, message []byte) {
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

		session.OnClose(func(session *link.Session, reason error) {
			println("client", session.Conn().RemoteAddr().String(), "close, ", reason)
		})
		session.Start()
	})

}
