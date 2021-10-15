package main

import (
	"JumboDB/jumbodb-core/src/server"
	"flag"
	"github.com/BurntSushi/toml"
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
	var conf tomlConfig
	if _, err := toml.DecodeFile(configPath, &conf); err != nil {
		log.Fatalf("Unable to parse JumboDB config in path %s", configPath)
	}

	server.StartListening(conf.Connection.Port)

}
