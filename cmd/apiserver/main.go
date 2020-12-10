package main

import (
	"github.com/MeguMan/geoTrain/internal/app/apiserver"
	"log"
)

func main() {
	if err := apiserver.Start(); err != nil {
		log.Fatal(err)
	}
}