package app

import (
	"os"

	"../eula"
)

func (a *Application) RunStage1() {
	a.LogInfo("Entering Run stage 1 - Initialisation Checks", nil, false)
	a.saveRunStage1()
}

func (a *Application) RunStage2() {
	a.LogInfo("Entering Run stage 2 - End User License Agreement checks", nil, false)
	if a.Config.GetBool("eula.agreed") {
		a.LogInfo("EULA has been agreed to.", nil, false)
		a.RunStage3()
	} else {
		a.LogInfo("EULA has not yest been accepted.", nil, false)
		answer := eula.DisplayEULA()
		if answer {
			a.Config.Set("eula.agreed", true)
			a.saveConfig()
			a.LogInfo("EULA acceptance state has been updated....Thankyou.", nil, false)
			a.LogInfo("Please rerun the application to continue.", nil, false)
			os.Exit(0)
		} else {
			a.LogInfo("The application cannot continue unless the EULA is accepted.", nil, false)
			os.Exit(0)
		}
	}
}

func (a *Application) RunStage3() {
	a.LogInfo("Entering Run stage 3 - Integrity Check", nil, false)
	if a.Status.eula == true {
		if a.Status.ucsCount > 1 {
			if a.Status.ucspmCount == 1 {
				a.LogInfo("All systems, config and checks completed successfully.", nil, false)
				a.RunStage4()
			} else {
				a.Log("The is no UCS Performance Manager system entered into the config file.", nil, false)
				os.Exit(0)
			}
		} else {
			a.Log("There is no UCS Systems entered into the config file.", nil, false)
			os.Exit(0)
		}
	} else {
		a.Log("The EULA needs to be agreed to before continuing.", nil, false)
		os.Exit(0)
	}
}

func (a *Application) RunStage4() {
	a.LogInfo("Entering Run stage 4 - UCS Performance Manager", nil, false)
	a.ucspmInit()
	a.ucspmInventory()
	a.saveRunStage4()
	a.RunStage5()
}

func (a *Application) RunStage5() {
	a.LogInfo("Entering Run stage 5 - UCS Manager Systems", nil, false)
	a.ucsInit()
	a.ucsInventory()
	a.saveRunStage5()
	a.RunStage6()
}

func (a *Application) RunStage6() {
	a.LogInfo("Entering Run stage 6 - UCS Performance Manager Reports", nil, false)
	a.ucspmProcessReports()
	a.saveRunStage6()
	a.RunStage7()
}

func (a *Application) RunStage7() {
	a.LogInfo("Entering Run stage 7 - Finialising and Saving reports", nil, false)
	a.saveRunStage7()
	//a.ucsInit()
	//a.ucsInventory()
}
