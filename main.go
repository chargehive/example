package main

import (
	"log"
	"os"
	"path/filepath"
)

var config = conf{}
var api = Api{}

func main() {
	webserver(&config)
}

func init() {
	if wd, err := os.Getwd(); err != nil {
		log.Fatal(err)
	} else if configFileName, err = filepath.Abs(filepath.Join(wd, configFileName)); err != nil {
		log.Fatal(err)
	} else if err = config.loadConfig(configFileName); err != nil {
		log.Fatal(err)
	}

	// populate missing fields with .chive.yaml if available
	if config.ProjectId == "" || config.ApiAccessToken == "" || config.ApiHost == "" {
		var chc chiveConf
		if err := chc.loadChiveConfig(); err != nil {
			log.Fatal("missing required fields in config and cannot load .chive.yaml")
		} else if err = config.populateConfigFromChive(&chc); err != nil {
			log.Fatal("missing required fields in config and cannot populate from .chive.yaml: " + err.Error())
		}
	}

	log.Printf("Using config %+v\n", config)
}
