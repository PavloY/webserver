package main

import (
	"WEB_SERVER/internal/app/api"
	"log"
	"flag"

	"github.com/BurntSushi/toml"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "path", "configs/api.toml", "path to config file in .toml format")
}

func main() {
	flag.Parse()
	log.Println("Starting ...")
	config := api.NewConfig()
	_, err := toml.DecodeFile(configPath, config)
	if err !=nil {
		log.Println("Can not find configs file. Using default values:", err)
	}

	server := api.New(config)
	
	if err := server.Start(); err !=nil {
		log.Fatal(err)
	}
}