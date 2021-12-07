package main

import (
	"JumboDB/jumboDB-core/src/config"
	"JumboDB/jumbodb-core/src/server"
	"flag"
	"log"
)

var (
	configPath string
)

//const logPath = "/some/location/example.log"

func init() {
	flag.StringVar(&configPath, "config", "config.toml", "JumboDB toml config path")
}

func main() {
	log.Println("Welcome to JumboDB")
	flag.Parse()
	config.Path = configPath
	server := server.NewHttpServer()
	server.StartListening()

}
