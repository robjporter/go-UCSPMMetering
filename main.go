package main

import (
	"fmt"

	"./flags"
)

func init() {
	flags.LoadConfig("./config.yaml")
}

func main() {
	ret := flags.ProcessCommandLineArguments()
	if ret == "RUN" {
		fmt.Println("HERE")
	}
}
