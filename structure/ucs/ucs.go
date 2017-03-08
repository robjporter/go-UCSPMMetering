package ucs

import (
	"fmt"

	"../app"
)

func PrintGreeting() {
	fmt.Println("UCS")
	app.Run("TEST")
	fmt.Println(app.App.ConfigFile)
	app.App.Run2()
}
