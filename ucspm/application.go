package ucspm

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"../flags"

	jmespath "github.com/jmespath/go-jmespath"
	functions "github.com/robjporter/go-functions"
	"github.com/robjporter/go-functions/as"
	"github.com/robjporter/go-functions/http"
	"github.com/robjporter/go-functions/viper"
)

var ()

func (a *Application) init() {
	a.routers = make(map[string]string)
	a.tidCount = 1
	a.host = ""
	a.username = ""
	a.password = ""
	a.outputFile = ""
	a.setupRouters()
}

func (a *Application) setupRouters() {
	a.routers["messaging"] = "MessagingRouter"
	a.routers["evconsole"] = "EventsRouter"
	a.routers["process"] = "ProcessRouter"
	a.routers["service"] = "ServiceRouter"
	a.routers["device"] = "DeviceRouter"
	a.routers["network"] = "NetworkRouter"
	a.routers["template"] = "TemplateRouter"
	a.routers["detailnav"] = "DetailNavRouter"
	a.routers["report"] = "ReportRouter"
	a.routers["mib"] = "MibRouter"
	a.routers["zenpack"] = "ZenPackRouter"
}

func (a *Application) CheckConfig(filename string) bool {
	a.init()
	configName := ""
	configExtension := ""
	configPath := ""

	splits := strings.Split(filepath.Base(filename), ".")
	if len(splits) == 2 {
		configName = splits[0]
		configExtension = splits[1]
	}
	configPath = filepath.Dir(filename)

	if functions.Exists(filename) {
		viper.SetConfigName(configName)
		viper.SetConfigType(configExtension)
		viper.AddConfigPath(configPath)
		a.configFile = filename
		return true
	} else {
		return false
	}
}

func (a *Application) Run() {
	a.loadConfig()
	a.setDefaults()
	a.start()
}

func (a *Application) setDefaults() {
	if viper.IsSet("output.file") {
		a.outputFile = viper.GetString("output.file")
	} else {
		a.outputFile = "output.csv"
	}
	if viper.IsSet("ucspm.url") {
		a.host = "https://" + viper.GetString("ucspm.url") + "/"
	}
	if viper.IsSet("ucspm.username") {
		a.username = viper.GetString("ucspm.username")
	}
	if viper.IsSet("ucspm.password") {
		a.password = flags.DecryptPassword(viper.GetString("ucspm.password"))
	}
}

func (a *Application) loadConfig() {
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}

func (a *Application) start() {
	err := errors.New("")
	a.devices, err = a.getDevices("device", "getDevices", `[{"uid": "/zport/dmd/Devices"}]`)
	if err == nil {
		a.markIgnoreUIDS()
		a.addHostsUnderVcenters()
		a.getUUIDForDevices()
		a.saveUUID(outputUUID())
	}
}

func (a *Application) saveUUID(json string) {
	f, err := os.Create(a.outputFile)
	if err == nil {
		_, err := f.Write([]byte(json))
		if err == nil {
			fmt.Println("\nFile saved successfully.")
		} else {
			fmt.Println("\nThere was an error saving the file: ", err)
		}
	}
	defer f.Close()
}

func (a *Application) outputUUID() string {
	jsonStr := `{"uuids": [`
	uuid := []string{}

	fmt.Println("Building UUID list")

	for i := 0; i < len(a.devices); i++ {
		if !a.devices[i].ignore {
			if a.devices[i].uuid != "" {
				uuid = append(uuid, a.devices[i].uuid)
			}
		}
	}

	fmt.Println("Removing duplicate UUID")
	uuid = removeDuplicates(uuid)

	for i := 0; i < len(uuid); i++ {
		jsonStr += `"` + uuid[i] + `",`
	}

	fmt.Println("Buidling JSON output")

	jsonStr = strings.TrimRight(jsonStr, ",")
	jsonStr += `]}`
	return jsonStr
}

func (a *Application) getUUIDForDevices() {
	fmt.Println("Getting UUID for all standalone servers")
	for i := 0; i < len(a.devices); i++ {
		if !a.devices[i].ignore {
			dev, err := getDeviceDetails(a.devices[i])
			if err == nil {
				tmp := device{}
				if dev != tmp {
					a.devices[i] = dev
				}
				if !a.devices[i].ignore {
					/*
						fmt.Println("NAME:> ", devices[i].name)
						fmt.Println("UID:> ", devices[i].uid)
						fmt.Println("UUID:> ", devices[i].uuid)
						fmt.Println("MODEL:> ", devices[i].model)
						fmt.Println("HYPERVISOR:> ", devices[i].hypervisor)
						fmt.Println("HYPERVISORNAME:> ", devices[i].hypervisorName)
						fmt.Println("HYPERVISORVERSION:> ", devices[i].hypervisorVersion)
						fmt.Println("HASHYPERVISOR:> ", devices[i].hasHypervisor)
						fmt.Println("UCSPMNAME:> ", devices[i].ucspmName)
						fmt.Println("=========================================")
					*/
				}
			} else {
				a.devices[i].ignore = true
			}
		}
	}
}

func (a *Application) generateUCSPMName(dev device) string {
	retName := "CPU_Utilization_-_-"
	if strings.Contains(dev.uid, "/zport/dmd/Devices/vSphere/devices") && strings.Contains(dev.uid, "datacenters") {
		retName += "vSphere-vCenter_-_" + dev.name
	} else {
		retName += "vSphere-" + dev.name + "_-_" + dev.hypervisorName
	}
	return retName
}

func (a *Application) getDeviceDetails(dev device) (device, error) {
	if dev.hasHypervisor == true {
		tmp, err := getHypervisorDeviceDetail(dev)
		if err == nil {
			return tmp, nil
		} else {
			return device{}, err
		}
	} else {
		if strings.Contains(dev.uid, "/zport/dmd/Devices/vSphere") {
			tmp, err := getStandaloneVsphereDeviceDetail(dev)
			if err == nil {
				return tmp, nil
			} else {
				return device{}, err
			}
		}
	}
	return device{}, errors.New("End of getDeviceDetails reached with no result")
}

func (a *Application) getStandaloneVsphereDeviceDetail(dev device) (device, error) {
	jsonStr := `{"action": "DeviceRouter","method": "getInfo","data":[{"uid": "` + dev.uid + `/datacenters/Datacenter_ha-datacenter/hosts/HostSystem_ha-host","keys": ["hardwareModel","hardwareUUID","hostname","name","hypervisorVersion","device"]}], "tid": ` + as.ToString(a.tidCount) + `}`
	url := a.host + strings.TrimLeft(dev.uid, "/") + "/device_router"
	headers := a.getHeaders()
	code, response, err := http.SendUnsecureHTTPSRequest(url, "POST", jsonStr, headers)
	a.tidCount++

	if err == nil {
		if code == 200 {
			if response != "" {

				var data2 interface{}
				json.Unmarshal([]byte(response), &data2)

				uuid, err := jmespath.Search("result.data.hardwareUUID", data2)
				model, err2 := jmespath.Search("result.data.hardwareModel", data2)
				hypname, err3 := jmespath.Search("result.data.name", data2)
				hypversion, err4 := jmespath.Search("result.data.hypervisorVersion", data2)
				name, err5 := jmespath.Search("result.data.device.name", data2)

				if err == nil && err2 == nil && err3 == nil && err4 == nil && err5 == nil {
					dev.hypervisor = true
					if as.ToString(uuid) == "" {
						dev.ignore = true
					} else {
						dev.name = as.ToString(name)
						dev.uuid = as.ToString(uuid)
						dev.model = as.ToString(model)
						dev.hypervisorName = as.ToString(hypname)
						dev.hypervisorVersion = as.ToString(hypversion)
						dev.ucspmName = generateUCSPMName(dev)
					}
					return dev, nil
				} else {
					return device{}, errors.New("Unknown hardware device")
				}
			}
		}
	}
	return device{}, errors.New("Unsuccessful connection")
}

func (a *Application) getHypervisorDeviceDetail(dev device) (device, error) {
	jsonStr := `{"action": "DeviceRouter","method": "getInfo","data":[{"uid": "` + dev.uid + `","keys": ["hardwareModel","id","hardwareUUID","uuid","hostname"]}],"tid":` + as.ToString(a.tidCount) + `}`
	url := a.host + "zport/dmd/device_router"
	headers := a.getHeaders()
	code, response, err := http.SendUnsecureHTTPSRequest(url, "POST", jsonStr, headers)
	a.tidCount++

	if err == nil {
		if code == 200 {
			if response != "" {

				var data2 interface{}
				json.Unmarshal([]byte(response), &data2)

				uuid, err := jmespath.Search("result.data.hardwareUUID", data2)
				model, err2 := jmespath.Search("result.data.hardwareModel", data2)
				name, err3 := jmespath.Search("result.data.hostname", data2)
				hypname, err4 := jmespath.Search("result.data.hostname", data2)
				fmt.Printf("Adding server: %s with UID: %s\n", as.ToString(name), dev.uid)

				if err == nil && err2 == nil && err3 == nil && err4 == nil {
					dev.name = as.ToString(name)
					dev.uuid = as.ToString(uuid)
					dev.model = as.ToString(model)
					dev.hypervisorName = as.ToString(hypname)
					dev.hypervisor = true
					dev.ucspmName = generateUCSPMName(dev)
					return dev, nil
				} else {
					return device{}, errors.New("Unknown hardware device")
				}
			}
		}
	}
	return device{}, errors.New("Unsuccessful connection")
}

func (a *Application) addHostsUnderVcenters() {
	fmt.Println("Adding servers under Hypervisors")
	for i := 0; i < len(a.devices); i++ {
		if a.devices[i].hypervisor {
			jsonStr := `{"action":"DeviceRouter","method":"getComponents","data":[{"uid":"/zport/dmd/Devices/vSphere/devices/vCenter","keys":["uid","id","title","name","hypervisorVersion","totalMemory","uuid"],"meta_type":"vSphereHostSystem","sort":"name","dir":"ASC"}],"tid":` + as.ToString(a.tidCount) + `}`
			url := a.host + "zport/dmd/Devices/vSphere/devices/vCenter/device_router"
			headers := a.getHeaders()
			code, response, err := http.SendUnsecureHTTPSRequest(url, "POST", jsonStr, headers)
			a.tidCount++
			a.devices[i].ignore = true

			if err == nil {
				if code == 200 {
					if response != "" {

						var data2 interface{}
						json.Unmarshal([]byte(response), &data2)

						tmp, err := jmespath.Search("result.totalCount", data2)
						if err == nil {
							for i := 0; i < int(as.ToInt(tmp)); i++ {
								version, err2 := jmespath.Search("result.data["+as.ToString(i)+"].hypervisorVersion", data2)
								uid, err4 := jmespath.Search("result.data["+as.ToString(i)+"].uid", data2)
								name, err5 := jmespath.Search("result.data["+as.ToString(i)+"].name", data2)
								fmt.Printf("Adding server: %s \n", as.ToString(name))
								if err2 == nil && err4 == nil && err5 == nil {
									var dev device
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

									a.devices = append(a.devices, dev)
								}
							}
						}
					}
				}
			}
		}
	}
}

func (a *Application) markIgnoreUIDS() {
	fmt.Println("Marking devices to be ignored")
	ignoreCountCompute := 0
	ignoreCountNetwork := 0
	ignoreCountStorage := 0
	for i := 0; i < len(a.devices); i++ {
		if strings.Contains(a.devices[i].uid, "zport/dmd/Devices/CiscoUCS/") {
			if a.DEBUG {
				fmt.Println("Ignoring UUID - Compute Device - " + a.devices[i].uid)
			}
			a.devices[i].ignore = true
			ignoreCountCompute++
		} else if strings.Contains(a.devices[i].uid, "/zport/dmd/Devices/Network/") {
			if a.DEBUG {
				fmt.Println("Ignoring UUID - Network Device - " + a.devices[i].uid)
			}
			a.devices[i].ignore = true
			ignoreCountNetwork++
		} else if strings.Contains(a.devices[i].uid, "/zport/dmd/Devices/Storage/") {
			if a.DEBUG {
				fmt.Println("Ignoring UUID - Storage Device - " + a.devices[i].uid)
			}
			a.devices[i].ignore = true
			ignoreCountStorage++
		}
	}
	totalCount := ignoreCountCompute + ignoreCountNetwork + ignoreCountStorage
	fmt.Println("Marked ", totalCount, " devices to be ignored. ", ignoreCountCompute, " Compute, ", ignoreCountNetwork, " Network and ", ignoreCountStorage, " Storage.")
}

func (a *Application) isVcenter(name string) bool {
	if strings.Contains(name, "VMware vCenter Server") {
		return true
	}
	return false
}

func (a *Application) getDevices(router string, method string, data string) ([]device, error) {
	fmt.Println("Getting all devices")
	devs := []device{}
	jsonStr := `{"action":"` + a.routers[router] + `","method":"` + method + `","data":` + data + `,"tid":` + as.ToString(a.tidCount) + `}`
	url := a.host + "zport/dmd/" + router + "_router"
	headers := a.getHeaders()
	code, response, err := http.SendUnsecureHTTPSRequest(url, "POST", jsonStr, headers)
	a.tidCount++

	if err == nil {
		if code == 200 {
			if response != "" {

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
									var tmp device
									tmp.uid = as.ToString(uid)
									tmp.ignore = false
									tmp.name = as.ToString(name)
									if err3 == nil {
										tmp.hypervisor = isVcenter(as.ToString(name))
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
	fmt.Println("Found ", len(devs), " devices to index.")
	return devs, nil
}

func (a *Application) getHeaders() map[string]string {
	headers := make(map[string]string)
	headers["Content-type"] = "application/json"
	headers["Accept-Charset"] = "utf-8"
	headers["Authorization"] = "Basic " + a.basicAuth(a.username, a.password)
	return headers
}

func (a *Application) basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func (a *Application) removeDuplicates(elements []string) []string {
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
