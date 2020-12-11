package main

import (
	"encoding/json"
	"github.com/MeguMan/geoTrain/internal/app/apiserver"
	"github.com/MeguMan/geoTrain/internal/app/memcache"
	"log"
	"os"
)

func main() {
	config := memcache.NewConfig()
	configFile, err := os.Open("configs/config.json")
	if err != nil {
		log.Fatal(err)
	}
	jsonParser := json.NewDecoder(configFile)
	if err := jsonParser.Decode(config); err != nil {
		log.Fatal(err)
	}

	if err := apiserver.Start(config); err != nil {
		log.Fatal(err)
	}
}