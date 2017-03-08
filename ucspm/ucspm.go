package ucspm

import (
	"fmt"

	"../app"
)

func PrintGreeting() {
	fmt.Println("UCSPM")
	app.Run("TEST")
	fmt.Println(app.App.ConfigFile)
	app.App.Run2()
}
