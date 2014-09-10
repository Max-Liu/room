package main

import (
	"fmt"
	"log"
	"net/http"
	"net/rpc"
	"room"
	"time"
)

var client *rpc.Client
var err error

func NewRoom(w http.ResponseWriter, r *http.Request) {
	var reply *[]byte
	var line []byte
	err := client.Call("RpcListener.CreateChatRoom", line, &reply)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, string(*reply))

}

func GetRoomNum(w http.ResponseWriter, r *http.Request) {

	var reply *[]byte
	var line []byte
	err := client.Call("RpcListener.GetRoomNum", line, &reply)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, string(*reply))

}

func EndRoom(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var reply *[]byte
	var line []byte
	line = []byte(r.Form["port"][0])
	err := client.Call("RpcListener.EndChatRoom", line, &reply)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("port %s has closed", string(line))
	fmt.Fprintf(w, string(*reply))
}
func GetRoomList(w http.ResponseWriter, r *http.Request) {
	var reply *[]byte
	var line []byte
	err := client.Call("RpcListener.GetRoomList", line, &reply)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, string(*reply))
}

var Env string = "development"

func main() {
	config := room.InitConfig(Env)
	address := config.GetDialAddress("chat")

	client, err = rpc.Dial("tcp", address)
	if err != nil {
		log.Fatal(err)
	}

	go Ping()

	http.HandleFunc("/new_room", NewRoom)
	http.HandleFunc("/delete", EndRoom)
	http.HandleFunc("/get_room_num", GetRoomNum)
	http.HandleFunc("/get_room_list", GetRoomList)
	err = http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
func Ping() {
	for {
		var reply *[]byte
		var line []byte
		log.Println("Sending heartbeat to server...")
		err := client.Call("RpcListener.Ping", line, &reply)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Recived heartbeat from server:", string(*reply), "Send next heartbeat in 10 second.")
		<-time.Tick(10 * time.Second)
	}
}
