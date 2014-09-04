package main

import (
	"fmt"
	"log"
	"net/http"
	"net/rpc"
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

func main() {
	client, err = rpc.Dial("tcp", "localhost:42586")
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/new_room", NewRoom)
	http.HandleFunc("/delete", EndRoom)
	http.HandleFunc("/get_room_num", GetRoomNum)
	http.HandleFunc("/get_room_List", GetRoomList)
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
