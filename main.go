package main

import (
	"fmt"
	"os"
	"runtime"

	"./flags"
	"./ucs"
)

var (
	ucsSystems ucs.Application
	configName = "./config.yaml"
)

func init() {
	flags.LoadConfig(configName)
	if !flags.EULACompliance() {
		fmt.Println("EULA NEEDS AGREEMENT")
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	ret := flags.ProcessCommandLineArguments()
	if ret == "RUN" {
		if flags.EULACompliance() {
			if !ucsSystems.CheckConfig(configName) {
				fmt.Println("The config file could not be found.  Please check and try again.")
				os.Exit(0)
			} else {
				ucsSystems.Run()
			}
		} else {
			fmt.Println("Unable to continue until EULA has been agreed.")
			os.Exit(0)
		}
	}
}
