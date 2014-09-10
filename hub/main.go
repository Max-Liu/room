package main

import "room"

var Env string = "development"

func main() {
	chatServerAddress := room.InitConfig(Env).GetDialAddress("chat")
	room.StartApiServer(chatServerAddress)
}
