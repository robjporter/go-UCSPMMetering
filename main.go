package main

import (
	"./app"
)

func main() {
	app.Core.ConfigFile = "./config.yaml"
	app.Core.LoadConfig()
	app.Core.Run()
}
