package main

import (
	"encoding/binary"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"room"

	"github.com/funny/link"
)

var (
	benchmark = flag.Bool("bench", false, "is for benchmark, will disable print")
)

func main() {
	flag.Parse()

	protocol := link.PacketN(2, binary.BigEndian)

	server, err := link.Listen("tcp", "127.0.0.1:10010", protocol)
	if err != nil {
		panic(err)
	}

	println("server start:", server.Listener().Addr().String())

	box := room.Box{}

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
