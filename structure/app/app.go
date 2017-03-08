package app

import (
	"fmt"
)

var (
	App Application
)

type Application struct {
	ConfigFile string
}

func init() {
}

func Run(filename string) {
	fmt.Println(filename)
	fmt.Println("INSIDE TESTME")
}

func (a *Application) Run2() {
	fmt.Println("APPLICATION.RUN")
}
