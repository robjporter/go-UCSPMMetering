package main

import (
	"./app"
)

func main() {
	app.Core.Debug = false
	app.Core.ConfigFile = "./config.yaml"
	app.Core.LoadConfig()
	app.Core.Run()
}
