package main

import (
	"./app"
	"./ucs"
	"./ucspm"
)

func main() {
	app.Run("config.yaml")
	app.App.ConfigFile = "CONFIGFILE"
	app.App.Run2()
	ucs.PrintGreeting()
	ucspm.PrintGreeting()
}
