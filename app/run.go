package app

import (
	"os"

	"../eula"
)

func (a *Application) RunStage2() {
	a.LogInfo("Entering Run stage 2", nil, false)
	if a.Config.GetBool("eula.agreed") {
		a.LogInfo("EULA has been agreed to", nil, false)
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
	a.LogInfo("Entering Run stage 3", nil, false)
}
