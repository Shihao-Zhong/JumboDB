package main

import (
	"flag"
	"github.com/BurntSushi/toml"
	"log"
	"../src/server"
)

var (
	configPath string
	verbose    bool
)

const logPath = "/some/location/example.log"

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

	StartListening(config.Connection.Port)

}
