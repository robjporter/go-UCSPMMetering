package app

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/robjporter/go-functions/as"
	"github.com/robjporter/go-functions/http"
	jmespath "github.com/robjporter/go-functions/jmespath"
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

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
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
									a.Log("UCS Performance Manager Device found", map[string]interface{}{"Name": name, "UID": uid}, true)
									tmp.uid = as.ToString(uid)
									tmp.ignore = false
									tmp.name = as.ToString(name)
									if err3 == nil {
										tmp.ishypervisor = a.isVcenter(as.ToString(name))
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

func (a *Application) getHeaders() map[string]string {
	headers := make(map[string]string)
	headers["Content-type"] = "application/json"
	headers["Accept-Charset"] = "utf-8"
	headers["Authorization"] = "Basic " + basicAuth(a.Config.GetString("ucspm.username"), a.DecryptPassword(a.Config.GetString("ucspm.password")))
	return headers
}

func (a *Application) isVcenter(name string) bool {
	if strings.Contains(name, "VMware vCenter Server") {
		a.LogInfo("Found a vCenter to index.", map[string]interface{}{"Name": name}, false)
		return true
	}
	return false
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

func (a *Application) ucspmInventory() {
	a.LogInfo("Preparing to run inventory on UCS Performance Manager.", nil, false)
	err := errors.New("")
	a.UCSPM.Devices, err = a.getDevices("device", "getDevices", `[{"uid": "/zport/dmd/Devices"}]`)
	if err == nil {
		a.ucspmAddHostsUnderVcenters()
		a.ucspmMarkDevicesToIgnore()
		a.ucspmGetUUIDForDevices()
		a.ucspmProcessDeviceDuplicates()
		a.ucspmSaveUUID(a.ucspmOutputUUID())
	}
}

func (a *Application) ucspmSaveUUID(json string) {
	filename := a.Config.GetString("input.file")
	f, err := os.Create(filename)
	if err == nil {
		_, err := f.Write([]byte(json))
		if err == nil {
			a.LogInfo("File has been saved successfully.", map[string]interface{}{"Filename": filename}, false)
		} else {
			a.LogInfo("There was a problem saving the file.", map[string]interface{}{"Error": err}, false)
		}
	}
	defer f.Close()
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

func (a *Application) ucspmRemoveDuplicates(elements []string) []string {
	// Use map to record duplicates as we find them.
	encountered := map[string]bool{}
	result := []string{}

	for v := range elements {
		if encountered[elements[v]] == true {
			// Do not add duplicate.
		} else {
			// Record this element as an encountered element.
			encountered[elements[v]] = true
			// Append to result slice.
			result = append(result, elements[v])
		}
	}
	// Return the new slice.
	return result
}

func (a *Application) ucspmGenerateUCSPMName(dev UCSPMDeviceInfo) string {
	retName := "CPU_Utilization_-_-"
	if strings.Contains(dev.uid, "/zport/dmd/Devices/vSphere/devices") && strings.Contains(dev.uid, "datacenters") {
		retName += "vSphere-vCenter_-_" + dev.name
	} else {
		retName += "vSphere-" + dev.name + "_-_" + dev.hypervisorName
	}
	return retName
}

func (a *Application) ucspmGetDeviceDetails(dev UCSPMDeviceInfo) (UCSPMDeviceInfo, error) {
	a.Log("Running deep inventory on interesting UID.", map[string]interface{}{"UID": dev.uid}, true)
	if dev.hasHypervisor == true {
		tmp, err := a.ucspmGetHypervisorDeviceDetail(dev)
		if err == nil {
			return tmp, nil
		} else {
			return UCSPMDeviceInfo{}, err
		}
	} else {
		if strings.Contains(dev.uid, "/zport/dmd/Devices/vSphere") {
			tmp, err := a.ucspmGetStandaloneVsphereDeviceDetail(dev)
			if err == nil {
				return tmp, nil
			} else {
				return UCSPMDeviceInfo{}, err
			}
		}
	}
	return UCSPMDeviceInfo{}, errors.New("End of getDeviceDetails reached with no result")
}

func (a *Application) ucspmGetUUIDForDevices() {
	a.LogInfo("Preparing to run deep inventory on interesting UID.", nil, false)
	for i := 0; i < len(a.UCSPM.Devices); i++ {
		if !a.UCSPM.Devices[i].ignore {
			dev, err := a.ucspmGetDeviceDetails(a.UCSPM.Devices[i])
			if err == nil {
				tmp := UCSPMDeviceInfo{}
				if dev != tmp {
					a.UCSPM.Devices[i] = dev
				}
			} else {
				a.UCSPM.Devices[i].ignore = true
			}
		}
	}
}

func (a *Application) ucspmGetStandaloneVsphereDeviceDetail(dev UCSPMDeviceInfo) (UCSPMDeviceInfo, error) {
	router := "DeviceRouter"
	method := "getInfo"
	data := `[{"uid": "` + dev.uid + `/datacenters/Datacenter_ha-datacenter/hosts/HostSystem_ha-host","keys": ["hardwareModel","hardwareUUID","hostname","name","hypervisorVersion","device"]}]`
	jsonStr := `{"action":"` + router + `","method":"` + method + `","data":` + data + `, "tid": ` + as.ToString(a.UCSPM.TidCount) + `}`
	url := a.makeUCSPMHostname() + strings.TrimLeft(dev.uid, "/") + "/device_router"
	headers := a.getHeaders()
	code, response, err := http.SendUnsecureHTTPSRequest(url, "POST", jsonStr, headers)
	a.UCSPM.TidCount++

	if err == nil {
		if code == 200 {
			if response != "" {
				a.LogInfo("Successfully received response from UCSPM.", map[string]interface{}{"Code": code}, true)

				var data2 interface{}
				json.Unmarshal([]byte(response), &data2)

				uuid, err := jmespath.Search("result.data.hardwareUUID", data2)
				model, err2 := jmespath.Search("result.data.hardwareModel", data2)
				hypname, err3 := jmespath.Search("result.data.name", data2)
				hypversion, err4 := jmespath.Search("result.data.hypervisorVersion", data2)
				name, err5 := jmespath.Search("result.data.device.name", data2)

				if err == nil && err2 == nil && err3 == nil && err4 == nil && err5 == nil {
					dev.ishypervisor = true
					if as.ToString(uuid) == "" {
						dev.ignore = true
					} else {
						dev.name = as.ToString(name)
						dev.uuid = as.ToString(uuid)
						dev.model = as.ToString(model)
						dev.hypervisorName = as.ToString(hypname)
						dev.hypervisorVersion = as.ToString(hypversion)
						dev.ucspmName = a.ucspmGenerateUCSPMName(dev)
						a.Log("UCS Performance Manager Device found", map[string]interface{}{"Name": name, "UUID": uuid}, true)

					}
					return dev, nil
				} else {
					return UCSPMDeviceInfo{}, errors.New("Unknown hardware device")
				}
			}
		}
	}
	a.Log("UCS Performance Manager Connection Error", map[string]interface{}{"Error": err, "UID": dev.uid}, true)
	return UCSPMDeviceInfo{}, errors.New("Unsuccessful connection")
}

func (a *Application) ucspmGetHypervisorDeviceDetail(dev UCSPMDeviceInfo) (UCSPMDeviceInfo, error) {
	router := "DeviceRouter"
	method := "getInfo"
	data := `[{"uid": "` + dev.uid + `","keys": ["hardwareModel","id","hardwareUUID","uuid","hostname"]}]`
	jsonStr := `{"action":"` + router + `","method":"` + method + `","data":` + data + `,"tid":` + as.ToString(a.UCSPM.TidCount) + `}`
	url := a.makeUCSPMHostname() + "zport/dmd/device_router"
	headers := a.getHeaders()
	code, response, err := http.SendUnsecureHTTPSRequest(url, "POST", jsonStr, headers)
	a.UCSPM.TidCount++

	if err == nil {
		if code == 200 {
			if response != "" {
				a.LogInfo("Successfully received response from UCSPM.", map[string]interface{}{"Code": code}, true)

				var data2 interface{}
				json.Unmarshal([]byte(response), &data2)

				uuid, err := jmespath.Search("result.data.hardwareUUID", data2)
				model, err2 := jmespath.Search("result.data.hardwareModel", data2)
				name, err3 := jmespath.Search("result.data.hostname", data2)
				hypname, err4 := jmespath.Search("result.data.hostname", data2)

				if err == nil && err2 == nil && err3 == nil && err4 == nil {
					dev.name = as.ToString(name)
					dev.uuid = as.ToString(uuid)
					dev.model = as.ToString(model)
					dev.hypervisorName = as.ToString(hypname)
					dev.ishypervisor = true
					dev.ucspmName = a.ucspmGenerateUCSPMName(dev)
					a.Log("UCS Performance Manager Device found", map[string]interface{}{"Name": name, "UUID": uuid}, true)
					return dev, nil
				} else {
					return UCSPMDeviceInfo{}, errors.New("Unknown hardware device")
				}
			}
		}
	}
	a.Log("UCS Performance Manager Connection Error", map[string]interface{}{"Error": err, "UID": dev.uid}, true)
	return UCSPMDeviceInfo{}, errors.New("Unsuccessful connection")
}

func (a *Application) ucspmAddHostsUnderVcenters() {
	router := "DeviceRouter"
	method := "getComponents"
	data := `[{"uid":"/zport/dmd/Devices/vSphere/devices/vCenter","keys":["uid","id","title","name","hypervisorVersion","totalMemory","uuid"],"meta_type":"vSphereHostSystem","sort":"name","dir":"ASC"}]`
	a.LogInfo("Preparing to inventory servers under discovered hypervisors.", map[string]interface{}{"Router": router, "Method": method, "Data": data}, false)
	count := 0
	for i := 0; i < len(a.UCSPM.Devices); i++ {
		if a.UCSPM.Devices[i].ishypervisor {
			jsonStr := `{"action":"` + router + `","method":"` + method + `","data":` + data + `,"tid":` + as.ToString(a.UCSPM.TidCount) + `}`
			url := a.makeUCSPMHostname() + "zport/dmd/Devices/vSphere/devices/vCenter/device_router"
			headers := a.getHeaders()

			code, response, err := http.SendUnsecureHTTPSRequest(url, "POST", jsonStr, headers)
			a.UCSPM.TidCount++
			a.UCSPM.Devices[i].ignore = true

			if err == nil {
				if code == 200 {
					if response != "" {
						a.LogInfo("Successfully received response from UCSPM.", map[string]interface{}{"Code": code}, true)

						var data2 interface{}
						json.Unmarshal([]byte(response), &data2)

						tmp, err := jmespath.Search("result.totalCount", data2)
						if err == nil {
							count = int(as.ToInt(tmp))
							for i := 0; i < int(as.ToInt(tmp)); i++ {
								version, err2 := jmespath.Search("result.data["+as.ToString(i)+"].hypervisorVersion", data2)
								uid, err4 := jmespath.Search("result.data["+as.ToString(i)+"].uid", data2)
								name, err5 := jmespath.Search("result.data["+as.ToString(i)+"].name", data2)
								a.Log("UCS Performance Manager Device found", map[string]interface{}{"Name": name, "UID": uid}, true)
								if err2 == nil && err4 == nil && err5 == nil {
									var dev UCSPMDeviceInfo
									dev.ignore = false
									dev.hasHypervisor = true

									if err2 == nil {
										dev.hypervisorVersion = as.ToString(version)
									}
									if err4 == nil {
										dev.uid = as.ToString(uid)
									}

									if err5 == nil {
										dev.name = as.ToString(name)
									}

									a.UCSPM.Devices = append(a.UCSPM.Devices, dev)
								}
							}
						}
					}
				}
			}
		}
	}
	a.LogInfo("Add all servers under hypervisors.", map[string]interface{}{"Servers": count}, false)
}

func (a *Application) ucspmMarkDevicesToIgnore() {
	a.LogInfo("Preparing to mark devices we are not interested in.", nil, false)
	count := 0
	for i := 0; i < len(a.UCSPM.Devices); i++ {
		if strings.Contains(a.UCSPM.Devices[i].uid, "zport/dmd/Devices/CiscoUCS/") {
			a.Log("Ignoring Cisco UCS Compute System", map[string]interface{}{"uid": a.UCSPM.Devices[i].uid}, true)
			a.UCSPM.Devices[i].ignore = true
			count++
		} else if strings.Contains(a.UCSPM.Devices[i].uid, "/zport/dmd/Devices/Network/") {
			a.Log("Ignoring Network device", map[string]interface{}{"uid": a.UCSPM.Devices[i].uid}, true)
			a.UCSPM.Devices[i].ignore = true
			count++
		} else if strings.Contains(a.UCSPM.Devices[i].uid, "/zport/dmd/Devices/Storage/") {
			a.Log("Ignoring storage device", map[string]interface{}{"uid": a.UCSPM.Devices[i].uid}, true)
			a.UCSPM.Devices[i].ignore = true
			count++
		}
	}
	a.LogInfo("Marked devices we are not going to index.", map[string]interface{}{"Ignored": count}, false)
}

func (a *Application) ucspmProcessDeviceDuplicates() {
	fmt.Println(a.ucspmGetNonIgnoredDevices())
	a.ucspmProcessDiscoveredDevices()
	fmt.Println(a.ucspmGetNonIgnoredDevices())
}

func (a *Application) ucspmProcessDiscoveredDevices() {

}

func (a *Application) ucspmGetNonIgnoredDevices() int {
	count := 0
	for i := 0; i < len(a.UCSPM.Devices); i++ {
		if !a.UCSPM.Devices[i].ignore {
			count++
		}
	}
	return count
}

func (a *Application) ucspmProcessReports() {
	a.LogInfo("Preparing to Process all data and request reports.", nil, false)
}
