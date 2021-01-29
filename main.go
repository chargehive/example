package main

import (
	"encoding/json"
	"github.com/chargehive/example/client"
	"github.com/chargehive/example/config"
	"github.com/chargehive/example/server"
	"log"
)

func main() {
	log.Println("Starting up example project...")
	config.Init()
	prettyConfig, _ := json.MarshalIndent(config.Get(), "", "\t")
	log.Printf("Using config:\n%s\n", prettyConfig)

	client.Init(config.Get())
	server.Start(config.Get())
}
