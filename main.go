package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"./flags"
	"./functions"
	"./ucs"
	"./ucspm"
)

var (
	ucsSystems     ucs.Application
	ucsPerformance ucspm.Application
	configName     = "./config.yaml"
)

func init() {
	flags.LoadConfig(configName)
	if !flags.EULACompliance() {
		answer := flags.DisplayEULA()
		if answer {
			fmt.Println("Thank you for accetping the End User License Agreemeent.  Please run the application again.")
		} else {
			fmt.Println("You will need to agree to the End User License Agreement before being able to use the application.")
		}
		os.Exit(0)
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	ret := flags.ProcessCommandLineArguments()
	splits := strings.Split(ret, "|")
	month := ""
	year := ""
	if len(splits) == 3 {
		if splits[0] == "RUN" {
			if splits[1] == "" {
				month = functions.CurrentMonthName()
			} else {
				tmp := functions.IsMonth(splits[1])
				if tmp != "" {
					month = tmp
				} else {
					month = functions.CurrentMonthName()
				}
			}
			if splits[2] == "" {
				year = functions.CurrentYear()
			} else {
				tmp := functions.IsYear(splits[2])
				if tmp != "" {
					year = tmp
				} else {
					year = functions.CurrentYear()
				}
			}
		}
	}
	if ret == "RUN" {
		if flags.EULACompliance() {
			if !ucsPerformance.CheckConfig(configName) {
				fmt.Println("The config file could not be found.  Please check and try again.")
				os.Exit(0)
			} else {
				ucsPerformance.Run(month, year)
				if !ucsSystems.CheckConfig(configName) {
					fmt.Println("The config file could not be found.  Please check and try again.")
					os.Exit(0)
				} else {
					ucsSystems.Run()
				}
			}
		} else {
			fmt.Println("Unable to continue until EULA has been agreed.")
			os.Exit(0)
		}
	}
}
