package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/rpc"
	"os/exec"
	"room"
)

var roomSysInfoList map[int]room.RoomSysInfo
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

func getRoomSysInfo(w http.ResponseWriter, r *http.Request) {
	cmd := exec.Command("go", "run", "main.go")
	stdout, _ := cmd.StdoutPipe()
	cmd.Start()
	line, _, _ := bufio.NewReader(stdout).ReadLine()

	roomSysInfo := room.RoomSysInfo{}
	json.Unmarshal(line, &roomSysInfo)

	roomSysInfoList[roomSysInfo.Pid] = roomSysInfo

	fmt.Fprintf(w, string(line))
}
