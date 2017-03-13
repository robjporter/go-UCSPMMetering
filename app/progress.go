package app

import (
	"strings"
	"time"

	functions "github.com/robjporter/go-functions"
	"github.com/robjporter/go-functions/as"
	"github.com/robjporter/go-functions/environment"
)

func (a *Application) saveRunStage1() {
	a.LogInfo("Saving data from Run Stage 1.", nil, false)

	jsonStr := `{"System": `
	jsonStr += "{"
	jsonStr += `"Time" : "` + as.ToString(time.Now()) + `",`
	jsonStr += `"isCompiled" : "` + as.ToString(environment.IsCompiled()) + `",`
	jsonStr += `"Compiler" : "` + environment.Compiler() + `",`
	jsonStr += `"CPU" : "` + as.ToString(environment.NumCPU()) + `",`
	jsonStr += `"Architecture" : "` + environment.GOARCH() + `",`
	jsonStr += `"OS" : "` + environment.GOOS() + `",`
	jsonStr += `"ROOT" : "` + environment.GOROOT() + `",`
	jsonStr += `"PATH" : "` + environment.GOPATH() + `"`
	jsonStr += `}}`

	a.saveFile("Stage1-SYS.json", jsonStr)
}

func (a *Application) saveRunStage2() {
	a.LogInfo("Saving data from Run Stage 2.", nil, false)
}

func (a *Application) saveRunStage3() {
	a.LogInfo("Saving data from Run Stage 3.", nil, false)
	err := functions.CopyFile(a.ConfigFile, "data/Stage3-Config.yaml")
	if err != nil {
		a.Log("Saving data from Run Stage 3 Failed.", map[string]interface{}{"Error": err}, false)
	} else {
		a.LogInfo("Saving data from Run Stage3 completed successfully.", nil, false)
	}
}

func (a *Application) saveRunStage4() {
	a.LogInfo("Saving data from Run Stage 4.", nil, false)
	//TODO: UCSPM Inventory

}

func (a *Application) saveRunStage5() {
	a.LogInfo("Saving data from Run Stage 5.", nil, false)

	jsonStr := `{"UCS": [`
	for i := 0; i < len(a.UCS.Systems); i++ {
		jsonStr += "{"
		jsonStr += `"Name" : "` + a.UCS.Systems[i].name + `",`
		jsonStr += `"IP" : "` + a.UCS.Systems[i].ip + `",`
		jsonStr += `"Version" : "` + a.UCS.Systems[i].version + `"`
		jsonStr += "},"
	}

	jsonStr = strings.TrimRight(jsonStr, ",")
	jsonStr += `]}`

	a.saveFile("Stage5-UCS.json", jsonStr)
}

func (a *Application) saveRunStage6() {
	a.LogInfo("Saving data from Run Stage 6.", nil, false)

	jsonStr := `{"Results": [`

	for i := 0; i < len(a.Results); i++ {

		jsonStr += "{"
		jsonStr += `"Name" : "` + a.Results[i].ucsName + `",`
		jsonStr += `"Description" : "` + a.Results[i].ucsDesc + `",`
		jsonStr += `"Model" : "` + a.Results[i].ucsModel + `",`
		jsonStr += `"Serial" : "` + a.Results[i].ucsSerial + `",`
		jsonStr += `"System" : "` + a.Results[i].ucsSystem + `",`
		jsonStr += `"Position" : "` + a.Results[i].ucsPosition + `",`
		jsonStr += `"DN" : "` + a.Results[i].ucsDN + `",`
		jsonStr += `"IsManaged" : "` + as.ToString(a.Results[i].isManaged) + `",`
		jsonStr += `"Name2" : "` + a.Results[i].ucspmName + `",`
		jsonStr += `"UID" : "` + a.Results[i].ucspmUID + `",`
		jsonStr += `"Key" : "` + a.Results[i].ucspmKey + `",`
		jsonStr += `"UUID" : "` + a.Results[i].ucspmUUID + `"`
		jsonStr += "},"
	}

	jsonStr = strings.TrimRight(jsonStr, ",")
	jsonStr += `]}`

	a.saveFile("Stage6-MergeResults.json", jsonStr)
}

func (a *Application) saveRunStage7() {
	a.LogInfo("Saving data from Run Stage 7.", nil, false)
	a.zipDataDir()
}
