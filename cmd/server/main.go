package main

import (
	"github.com/GabrielNSD/Go-music-link-api/configs"
)

func main() {
	config, _ := configs.LoadConfig(".")
	println(config.DBDriver)
}
