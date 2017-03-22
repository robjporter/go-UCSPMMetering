package main

import (
	"./app"
)

func main() {
	app.Core.Start()
	app.Core.Version = "0.5.3b"
	app.Core.ConfigFile = "./config.yaml"
	app.Core.LoadConfig()
	app.Core.Run()
	app.Core.Finish()
}
