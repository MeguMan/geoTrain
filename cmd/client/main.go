package main

import (
	"github.com/MeguMan/geoTrain/internal/app/client"
)

func main() {
	c := client.NewClient("http://localhost:8080")
	c.Login("supersecretpassword")
	c.Set("mykesasmdy2", "myvaluLKABDe", 0)
	c.Get("mykesasmdy2")
	c.HSet("myhash22", "myfield3", "myvalue4")
	c.HGet("myhash22", "myfield3")
	c.Delete("mykesasmdy2")
	c.Save()
}
