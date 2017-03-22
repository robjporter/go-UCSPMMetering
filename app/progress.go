package app

import (
	"os"
	"strings"
	"time"

	functions "github.com/robjporter/go-functions"
	"github.com/robjporter/go-functions/as"
	"github.com/robjporter/go-functions/environment"
)

func (a *Application) saveRunStage1() {
	if a.Action != "CLEAN" {
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
		jsonStr += `"PATH" : "` + environment.GOPATH() + `",`
		jsonStr += `"APPVERSION: " : "` + a.Version + `",`
		jsonStr += `"Args" : "` + as.ToString(os.Args) + `"`
		jsonStr += `}}`

		a.saveFile("Stage1-SYS.json", jsonStr)
	}
}

func (a *Application) saveRunStage2() {
	a.LogInfo("Saving data from Run Stage 2.", nil, false)
	//TODO:
}

func (a *Application) saveRunStage3() {
	a.LogInfo("Saving data from Run Stage 3.", nil, false)
	err := functions.CopyFile(a.ConfigFile, "data/Stage3-Config.yaml")
	if err != nil {
		a.Log("Saving data from Run Stage 3 Failed.", map[string]interface{}{"Error": err}, false)
	} else {
		a.LogInfo("Saving data from Run Stage 3 completed successfully.", nil, false)
	}
}

func (a *Application) saveRunStage4() {
	a.LogInfo("Saving data from Run Stage 4.", nil, false)
	//TODO: UCSPM Inventory\
	a.ucspmSaveUUID(a.ucspmOutputUUID())

}

func (a *Application) saveRunStage5() {
	a.LogInfo("Saving data from Run Stage 5.", nil, false)
	a.LogInfo("Saving all UCS System info.", nil, false)

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

	a.saveFile("Stage5-UCSSystems.json", jsonStr)

	a.saveUUIDS()
	a.saveIgnored()
}

func (a *Application) saveUUIDS() {
	a.LogInfo("Saving all server node info.", nil, false)
	jsonStr := `{"Servers": [`

	for i := 0; i < len(a.UCS.Matches); i++ {
		jsonStr += "{"
		jsonStr += `"UUID" : "` + a.UCS.Matches[i].serveruuid + `",`
		jsonStr += `"OUUID" : "` + a.UCS.Matches[i].serverouuid + `",`
		jsonStr += `"DN" : "` + a.UCS.Matches[i].serverdn + `",`
		jsonStr += `"DESCRIPTION" : "` + a.UCS.Matches[i].serverdescr + `",`
		jsonStr += `"POSITION" : "` + a.UCS.Matches[i].serverposition + `",`
		jsonStr += `"NAME" : "` + a.UCS.Matches[i].servername + `",`
		jsonStr += `"PID" : "` + a.UCS.Matches[i].serverpid + `",`
		jsonStr += `"MODEL" : "` + a.UCS.Matches[i].servermodel + `",`
		jsonStr += `"SERIAL" : "` + a.UCS.Matches[i].serverserial + `",`
		jsonStr += `"DOMAINNAME" : "` + a.UCS.Matches[i].ucsname + `",`
		jsonStr += `"DOMAINVERSION" : "` + a.UCS.Matches[i].ucsversion + `",`
		jsonStr += `"DOMAINURL" : "` + a.UCS.Matches[i].ucsip + `"`
		jsonStr += "},"
	}

	jsonStr = strings.TrimRight(jsonStr, ",")
	jsonStr += `]}`

	a.saveFile("Stage5-UCSServers.json", jsonStr)

}

func (a *Application) saveIgnored() {
	a.LogInfo("Saving all ignored device info.", nil, false)
	jsonStr := `{"Devices": [`

	for i := 0; i < len(a.UCSPM.Devices); i++ {
		if a.UCSPM.Devices[i].ignore {
			jsonStr += "{"
			jsonStr += `"hasHypervisor":"` + as.ToString(a.UCSPM.Devices[i].hasHypervisor) + `",`
			jsonStr += `"hypervisorName":"` + as.ToString(a.UCSPM.Devices[i].hypervisorName) + `",`
			jsonStr += `"hypervisorVersion":"` + as.ToString(a.UCSPM.Devices[i].hypervisorVersion) + `",`
			jsonStr += `"ignore":"` + as.ToString(a.UCSPM.Devices[i].ignore) + `",`
			jsonStr += `"isHypervisor":"` + as.ToString(a.UCSPM.Devices[i].ishypervisor) + `",`
			jsonStr += `"model":"` + as.ToString(a.UCSPM.Devices[i].model) + `",`
			jsonStr += `"name":"` + as.ToString(a.UCSPM.Devices[i].name) + `",`
			jsonStr += `"ucspmName":"` + as.ToString(a.UCSPM.Devices[i].ucspmName) + `",`
			jsonStr += `"uid":"` + as.ToString(a.UCSPM.Devices[i].uid) + `",`
			jsonStr += `"uuid":"` + as.ToString(a.UCSPM.Devices[i].uuid) + `"`
			jsonStr += "},"
		}
	}

	jsonStr = strings.TrimRight(jsonStr, ",")
	jsonStr += `]}`

	a.saveFile("Stage5-IgnoredDevices.json", jsonStr)
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

	a.saveFile("Stage6-MergedResults.json", jsonStr)

	a.LogInfo("Successfully matched UUIDs.", map[string]interface{}{"Discovered": len(a.UCS.UUID), "Matched": len(a.UCS.Matched)}, true)
	a.saveMatchedUUID()
	if len(a.UCS.Matched) < len(a.UCS.UUID) {
		a.LogInfo("There were some unmatched UUID's.", map[string]interface{}{"Unmatched": a.UCS.Unmatched}, true)
		a.saveUnmatchedUUID()
	}
}

func (a *Application) saveRunStage7() {
	a.LogInfo("Saving data from Run Stage 7.", nil, false)
	a.exportHTTPCommands()
	a.zipDataDir()
}

func (a *Application) exportHTTPCommands() {
	a.LogInfo("Exporting all HTTP requests and responses.", nil, false)

	jsonStr := `{"Results": [`

	for i := 0; i < len(a.Commands); i++ {
		jsonStr += "{"
		jsonStr += `"Request" : {`
		jsonStr += `"URL" : "` + a.Commands[i].RequestURL + `",`
		jsonStr += `"Headers" : "` + as.ToString(a.Commands[i].RequestHeaders) + `",`
		jsonStr += `"Body" : "` + a.Commands[i].RequestBody + `"`
		jsonStr += `},`
		jsonStr += `"Response" : {`
		jsonStr += `"Code" : "` + as.ToString(a.Commands[i].ResponseCode) + `",`
		jsonStr += `"Body" : "` + a.Commands[i].ResponseBody + `",`
		jsonStr += `"Error" : "` + as.ToString(a.Commands[i].ResponseError) + `"`
		jsonStr += `}`
		jsonStr += "},"
	}

	jsonStr = strings.TrimRight(jsonStr, ",")
	jsonStr += `]}`

	a.saveFile("Stage7-HTTPRequests.json", jsonStr)
}

func (a *Application) saveMatchedUUID() {
	a.LogInfo("Saving unmatched UUID.", map[string]interface{}{"Unmatched": len(a.UCS.Unmatched)}, false)
	found := false
	jsonStr := `{"UUIDS": [`
	for i := 0; i < len(a.UCS.Matched); i++ {
		for j := len(a.UCSPM.Devices) - 1; j > -1; j-- {
			if a.UCSPM.Devices[j].uuid == a.UCS.Matched[i].serveruuid && !found {
				jsonStr += "{"
				jsonStr += `"hasHypervisor":"` + as.ToString(a.UCSPM.Devices[j].hasHypervisor) + `",`
				jsonStr += `"hypervisorName":"` + as.ToString(a.UCSPM.Devices[j].hypervisorName) + `",`
				jsonStr += `"hypervisorVersion":"` + as.ToString(a.UCSPM.Devices[j].hypervisorVersion) + `",`
				jsonStr += `"ignore":"` + as.ToString(a.UCSPM.Devices[j].ignore) + `",`
				jsonStr += `"isHypervisor":"` + as.ToString(a.UCSPM.Devices[j].ishypervisor) + `",`
				jsonStr += `"model":"` + as.ToString(a.UCSPM.Devices[j].model) + `",`
				jsonStr += `"name":"` + as.ToString(a.UCSPM.Devices[j].name) + `",`
				jsonStr += `"ucspmName":"` + as.ToString(a.UCSPM.Devices[j].ucspmName) + `",`
				jsonStr += `"uid":"` + as.ToString(a.UCSPM.Devices[j].uid) + `",`
				jsonStr += `"uuid":"` + as.ToString(a.UCSPM.Devices[j].uuid) + `"`
				jsonStr += "},"
				found = true
			}
		}
		found = false
	}
	jsonStr = strings.TrimRight(jsonStr, ",")
	jsonStr += `]}`

	a.saveFile("Stage6-MatchedUUID.json", jsonStr)
}

func (a *Application) saveUnmatchedUUID() {
	a.LogInfo("Saving unmatched UUID.", map[string]interface{}{"Unmatched": len(a.UCS.Unmatched)}, false)
	found := false
	jsonStr := `{"UUIDS": [`
	for i := 0; i < len(a.UCS.Unmatched); i++ {
		for j := len(a.UCSPM.Devices) - 1; j > -1; j-- {
			if a.UCSPM.Devices[j].uuid == a.UCS.Unmatched[i] && !found {
				jsonStr += "{"
				jsonStr += `"hasHypervisor":"` + as.ToString(a.UCSPM.Devices[j].hasHypervisor) + `",`
				jsonStr += `"hypervisorName":"` + as.ToString(a.UCSPM.Devices[j].hypervisorName) + `",`
				jsonStr += `"hypervisorVersion":"` + as.ToString(a.UCSPM.Devices[j].hypervisorVersion) + `",`
				jsonStr += `"ignore":"` + as.ToString(a.UCSPM.Devices[j].ignore) + `",`
				jsonStr += `"isHypervisor":"` + as.ToString(a.UCSPM.Devices[j].ishypervisor) + `",`
				jsonStr += `"model":"` + as.ToString(a.UCSPM.Devices[j].model) + `",`
				jsonStr += `"name":"` + as.ToString(a.UCSPM.Devices[j].name) + `",`
				jsonStr += `"ucspmName":"` + as.ToString(a.UCSPM.Devices[j].ucspmName) + `",`
				jsonStr += `"uid":"` + as.ToString(a.UCSPM.Devices[j].uid) + `",`
				jsonStr += `"uuid":"` + as.ToString(a.UCSPM.Devices[j].uuid) + `"`
				jsonStr += "},"
				found = true
			}
		}
	}
	jsonStr = strings.TrimRight(jsonStr, ",")
	jsonStr += `]}`

	a.saveFile("Stage6-UnmatchedUUID.json", jsonStr)
}

func (a *Application) ucspmSaveUUID(json string) {
	a.saveFile("Stage4-DiscoveredUUID.json", json)
}

func (a *Application) ucspmOutputUUID() string {
	jsonStr := `{"uuids": [`
	uuid := []string{}

	a.LogInfo("Building identified UUID list.", nil, false)

	for i := 0; i < len(a.UCSPM.Devices); i++ {
		if !a.UCSPM.Devices[i].ignore {
			if a.UCSPM.Devices[i].uuid != "" {
				uuid = append(uuid, a.UCSPM.Devices[i].uuid)
			} else {
				a.UCSPM.Devices[i].ignore = true
			}
		}
	}

	a.LogInfo("Removing duplicates from UUID list.", nil, false)
	uuid = a.ucspmRemoveDuplicates(uuid)
	a.UCSPM.ProcessedUUID = uuid
	a.LogInfo("Identified unique UUID list.", map[string]interface{}{"UUID": len(uuid)}, false)

	for i := 0; i < len(uuid); i++ {
		jsonStr += `"` + uuid[i] + `",`
	}

	a.LogInfo("Building JSON output string", nil, false)

	jsonStr = strings.TrimRight(jsonStr, ",")
	jsonStr += `]}`
	return jsonStr
}
