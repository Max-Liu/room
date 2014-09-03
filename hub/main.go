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

func main() {
	client, err = rpc.Dial("tcp", "localhost:42586")
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/new_room", NewRoom)
	http.HandleFunc("/delete", EndRoom)
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
