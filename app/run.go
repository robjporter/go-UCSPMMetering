package app

import (
	"os"

	"../eula"
)

func (a *Application) RunStage2() {
	a.LogInfo("Entering Run stage 2", nil, false)
	if a.Config.GetBool("eula.agreed") {

		a.LogInfo("EULA has been agreed to", nil, false)
	} else {
		a.LogInfo("EULA has not yest been accepted.", nil, false)
		answer := eula.DisplayEULA()
		if answer {
			a.Config.Set("eula.agreed", true)
			a.saveConfig()
			a.LogInfo("EULA acceptance state has been updated....Thankyou.", nil, false)
		} else {
			a.LogInfo("The application cannot continue unless the EULA is accepted.", nil, false)
			os.Exit(0)
		}
	}
	a.RunStage3()
}

func (a *Application) RunStage3() {
	a.LogInfo("Entering Run stage 3", nil, false)
}
