package app

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	jmespath "github.com/jmespath/go-jmespath"
	"github.com/robjporter/go-functions/as"
	"github.com/robjporter/go-functions/http"
)

func (a *Application) ucspmInit() {
	a.UCSPM.Routers = make(map[string]string)
	a.UCSPM.Routers["messaging"] = "MessagingRouter"
	a.UCSPM.Routers["evconsole"] = "EventsRouter"
	a.UCSPM.Routers["process"] = "ProcessRouter"
	a.UCSPM.Routers["service"] = "ServiceRouter"
	a.UCSPM.Routers["device"] = "DeviceRouter"
	a.UCSPM.Routers["network"] = "NetworkRouter"
	a.UCSPM.Routers["template"] = "TemplateRouter"
	a.UCSPM.Routers["detailnav"] = "DetailNavRouter"
	a.UCSPM.Routers["report"] = "ReportRouter"
	a.UCSPM.Routers["mib"] = "MibRouter"
	a.UCSPM.Routers["zenpack"] = "ZenPackRouter"

	a.UCSPM.TidCount = 1
}

func (a *Application) ucspmInventory() {
	err := errors.New("")
	a.UCSPM.Devices, err = a.getDevices("device", "getDevices", `[{"uid": "/zport/dmd/Devices"}]`)
	if err == nil {
		fmt.Println("CONTINUING")
	}
}

func (a *Application) makeUCSPMHostname() string {
	tmp := a.Config.GetString("ucspm.url")
	if strings.Contains(tmp, "http") {
		return tmp
	} else {
		tmp = "https://" + tmp + "/"
	}
	return tmp
}

func (a *Application) getDevices(router string, method string, data string) ([]UCSPMDeviceInfo, error) {
	a.LogInfo("Getting all UCS Performance Manager Devices", map[string]interface{}{"Router": router, "Method": method, "Data": data}, false)
	devs := []UCSPMDeviceInfo{}
	jsonStr := `{"action":"` + a.UCSPM.Routers[router] + `","method":"` + method + `","data":` + data + `,"tid":` + as.ToString(a.UCSPM.TidCount) + `}`
	url := a.makeUCSPMHostname() + "zport/dmd/" + router + "_router"
	headers := a.getHeaders()
	code, response, err := http.SendUnsecureHTTPSRequest(url, "POST", jsonStr, headers)
	a.UCSPM.TidCount++

	if err == nil {
		if code == 200 {
			if response != "" {

				a.LogInfo("Successfully received response from UCSPM.", map[string]interface{}{"Code": code}, true)
				var data2 interface{}
				json.Unmarshal([]byte(response), &data2)

				tmp, err := jmespath.Search("result.totalCount", data2)
				if err == nil {
					for i := 0; i < int(as.ToInt(tmp)); i++ {
						uid, err2 := jmespath.Search("result.devices["+as.ToString(i)+"].uid", data2)
						name, err3 := jmespath.Search("result.devices["+as.ToString(i)+"].osModel.name", data2)
						control, err4 := jmespath.Search("result.devices["+as.ToString(i)+"].pythonClass", data2)
						if err2 == nil {
							if err4 == nil {
								if !strings.Contains(as.ToString(control), "ZenPacks.zenoss.ControlCenter.ControlCenter") {
									var tmp UCSPMDeviceInfo
									tmp.uid = as.ToString(uid)
									tmp.ignore = false
									tmp.name = as.ToString(name)
									if err3 == nil {
										tmp.hypervisor = a.isVcenter(as.ToString(name))
									}
									devs = append(devs, tmp)
								}
							}
						} else {
							return nil, err2
						}
					}
				} else {
					return nil, err
				}
			}
		}
	} else {
		return nil, err
	}
	a.LogInfo("UCSPM responded with devices to index.", map[string]interface{}{"Devices": len(devs)}, false)
	return devs, nil
}

func (a *Application) isVcenter(name string) bool {
	if strings.Contains(name, "VMware vCenter Server") {
		a.LogInfo("Found a vCenter to index.", map[string]interface{}{"Name": name}, false)
		return true
	}
	return false
}

func (a *Application) getHeaders() map[string]string {
	headers := make(map[string]string)
	headers["Content-type"] = "application/json"
	headers["Accept-Charset"] = "utf-8"
	headers["Authorization"] = "Basic " + basicAuth(a.Config.GetString("ucspm.username"), a.DecryptPassword(a.Config.GetString("ucspm.password")))
	return headers
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
