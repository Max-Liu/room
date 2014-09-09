package room

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Server struct {
	Id   string
	Host string
	Port int
}

type ChatServer struct {
	Server
}

type AuthServer struct {
	Server
}
type ServerConfig struct {
	Chat []ChatServer
	Auth []AuthServer
}

type Config struct {
	Development ServerConfig
	Production  ServerConfig
}

func InitConfig(Env string) *ServerConfig {
	configByte, err := ioutil.ReadFile("../config.json")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(configByte))

	newConfig := new(Config)
	err = json.Unmarshal(configByte, newConfig)
	if err != nil {
		log.Fatal(err)
	}
	var serverConfig *ServerConfig

	switch Env {
	case "development":
		{
			serverConfig = &newConfig.Development
		}
	case "production":
		{
			serverConfig = &newConfig.Production

		}
	default:
		{
			log.Fatal("wrong Env config string")
		}

	}
	return serverConfig

}
