package app

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"../functions"

	"github.com/robjporter/go-functions/as"
	"github.com/robjporter/go-functions/http"
	"github.com/robjporter/go-functions/jmespath"
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
										fmt.Println("ISVCENTER:> ", as.ToString(name))
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
	if strings.Contains(strings.ToLower(name), "vmware vcenter server") {
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
		a.ucspmSaveUUID(a.ucspmOutputUUID())
	}
}

func (a *Application) ucspmSaveUUID(json string) {
	filename := a.Config.GetString("output.matched")
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
	//TODO: THE DATA VALUE CANNOT BE STATICALLY SET, IT NEEDS TO BE DYNAMIC /devices/vCenter
	// Each UCSPM object, shoould have vcenter UID
	router := "DeviceRouter"
	method := "getComponents"
	count := 0
	for i := 0; i < len(a.UCSPM.Devices); i++ {
		if a.UCSPM.Devices[i].ishypervisor {
			data := `[{"uid":"` + a.UCSPM.Devices[i].uid + `","keys":["uid","id","title","name","hypervisorVersion","totalMemory","uuid"],"meta_type":"vSphereHostSystem","sort":"name","dir":"ASC"}]`
			jsonStr := `{"action":"` + router + `","method":"` + method + `","data":` + data + `,"tid":` + as.ToString(a.UCSPM.TidCount) + `}`
			url := a.makeUCSPMHostname() + strings.TrimLeft(a.UCSPM.Devices[i].uid, "/") + "/device_router"
			headers := a.getHeaders()
			a.LogInfo("Preparing to inventory servers under discovered hypervisors.", map[string]interface{}{"Router": router, "Method": method, "Data": data, "URL": url}, false)

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
	a.LogInfo("Removing duplicates recevied from UCS Performance Manager.", nil, false)
	originalCount := a.ucspmGetNonIgnoredDevices()
	a.ucspmProcessDiscoveredDevices()
	for i := 0; i < len(a.UCSPM.Devices); i++ {
		if !a.UCSPM.Devices[i].ignore {
			var tmp CombinedResults
			tmp.ucspmName = a.UCSPM.Devices[i].name
			tmp.ucspmUID = a.UCSPM.Devices[i].uid
			tmp.ucspmUUID = a.UCSPM.Devices[i].uuid
			tmp2 := a.ucsGetUCSSystem(tmp.ucspmUUID)
			tmp.ucsDN = tmp2.serverdn
			tmp.ucsDesc = tmp2.serverdescr
			tmp.ucsModel = tmp2.servermodel
			tmp.ucsName = tmp2.servername
			tmp.ucsPosition = tmp2.serverposition
			tmp.ucsSerial = tmp2.serverserial
			tmp.ucsSystem = tmp2.ucsname
			tmp.isManaged = a.UCSPM.Devices[i].hasHypervisor
			a.Results = append(a.Results, tmp)
		}
	}
	updatedCount := a.ucspmGetNonIgnoredDevices()
	a.LogInfo("Removed duplicates recevied from UCS Performance Manager.", map[string]interface{}{"OriginalUUID": originalCount, "CleanUUID": updatedCount}, false)
}

func (a *Application) ucspmProcessDiscoveredDevices() {
	var matched []string
	for i := len(a.UCSPM.Devices) - 1; i > -1; i-- {
		if !inStringSlice(matched, a.UCSPM.Devices[i].uuid) {
			matched = append(matched, a.UCSPM.Devices[i].uuid)
		} else {
			a.UCSPM.Devices[i].ignore = true
		}
	}
}

func inStringSlice(slice []string, needle string) bool {
	for i := 0; i < len(slice); i++ {
		if strings.TrimSpace(needle) == strings.TrimSpace(slice[i]) {
			return true
		}
	}
	return false
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
	a.LogInfo("Preparing to Process all data and request reports.", map[string]interface{}{"Requests": len(a.Results)}, false)

	a.ucspmProcessDeviceDuplicates()
	for i := 0; i < len(a.Results); i++ {
		if a.Results[i].isManaged {
			a.Results[i].ucspmKey = createUCSPMKey(a.Results[i].ucspmUID)
			a.ucspmGetManagedReport(a.Results[i])
		} else {
			a.ucspmGetUnmanagedReport(a.Results[i])
		}
	}
}

func (a *Application) saveRunStage4() {
	//TODO:
	fmt.Println("RUN STAGE 4")
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

func (a *Application) ucspmGetManagedReport(sys CombinedResults) {
	a.LogInfo("Preparing to request all UCS Performance Manager reports, for managed devices.", nil, false)
	start := functions.GetTimestampStartOfMonth(a.Report.Month, int(as.ToInt(a.Report.Year)))
	end := functions.GetTimestampEndOfMonth(a.Report.Month, int(as.ToInt(a.Report.Year)))
	jsonStr := `
			{
	"start": ` + as.ToString(start) + `,
	"end": ` + as.ToString(end) + `,
	"series": true,
	"downsample": "1h-avg",
	"tags": {},
	"returnset": "EXACT",
	"metrics": [{
		"metric": "vCenter/cpuUsage_cpuUsage",
		"rate": false,
		"rateOptions": {},
		"aggregator": "avg",
		"tags": {
			"key": ["` + sys.ucspmKey + `"]
		},
		"name": "Usage-raw",
		"emit": false
	}, {
		"name": "Usage",
		"expression": "rpn:Usage-raw,100,/"
	}]
}
		`
	url := a.makeUCSPMHostname() + "api/performance/query/"
	headers := a.getHeaders()
	a.LogInfo("Requesting report.", map[string]interface{}{"ReportStart": start, "ReportEnd": end, "UID": sys.ucspmUID, "Key": sys.ucspmKey, "URL": url}, false)

	code, response, err := http.SendUnsecureHTTPSRequest(url, "POST", jsonStr, headers)

	if err == nil {
		if code == 200 {
			if response != "" {
				a.LogInfo("Successfully received response from UCSPM.", map[string]interface{}{"Code": code}, true)

				var data2 interface{}
				json.Unmarshal([]byte(response), &data2)

				tmp, err := jmespath.Search("results[0].datapoints", data2)
				if err == nil {
					tmp2 := as.ToSlice(tmp)
					a.LogInfo("Received Datapoints to process.", map[string]interface{}{"Datapoints": len(tmp2)}, true)
					a.processReport(sys, tmp2)
				}
			}
		}
	}
}

func (a *Application) processReport(sys CombinedResults, data []interface{}) {
	m := make(map[int]ReportData)
	for i := 0; i < len(data); i++ {
		tmp := as.ToStringMap(data[i])
		ttmp := as.ToInt(strconv.FormatFloat(as.ToFloat(tmp["timestamp"]), 'f', 0, 64))
		var temp ReportData
		temp.timestamp = int(ttmp)
		temp.value = as.ToFloat(tmp["value"])
		m[i] = temp
	}
	s := make(dataSlice, 0, len(m))
	for _, d := range m {
		s = append(s, d)
	}
	sort.Sort(s)
	a.outputProcessedReport(sys, s)
}

func (a *Application) outputProcessedReport(sys CombinedResults, data dataSlice) {
	name := ""
	if sys.ucspmName != "" {
		name = sys.ucspmName
	} else if sys.ucsName != "" {
		name = sys.ucsName
	} else {

		rand.Seed(time.Now().UTC().UnixNano())
		name = "server" + as.ToString(rand.Intn(1000-9999))
	}

	filename := name + "-" + sys.ucsSerial + "-" + a.Report.Month + "-" + as.ToString(a.Report.Year) + "-" + as.ToString(time.Now().Unix()) + ".csv"

	csv := "timestamp,value\n"
	for _, d := range data {
		csv += as.ToString(d.timestamp) + "," + as.ToString(d.value) + "\n"
	}
	a.saveFile(filename, csv)
}

func (a *Application) ucspmGetUnmanagedReport(sys CombinedResults) {

}

func createUCSPMKey(uid string) string {
	name := "/zport/dmd/Devices/vSphere/d"
	newUID := ""
	if strings.HasPrefix(uid, name) {
		newUID = "D" + uid[len(name):len(uid)]
	}
	splits := strings.Split(newUID, "/")
	if len(splits) > 2 {
		splits[1] = "vCenter"
	}
	newerUID := ""
	for i := 0; i < len(splits); i++ {
		newerUID += splits[i] + "/"
	}
	newerUID = strings.TrimRight(newerUID, "/")
	return newerUID
}
