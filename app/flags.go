package app

import (
	"strings"

	"github.com/robjporter/go-functions/as"
)

func (a *Application) addUCS(ip, username, password string) bool {
	if ip != "" {
		if username != "" {
			if password != "" {
				tmp := UCSSystemInfo{}
				tmp.ip = ip
				tmp.username = username
				tmp.password = a.EncryptPassword(password)
				a.UCSSystems = append(a.UCSSystems, tmp)
				return true
			} else {
				a.Log("The password for the UCS System cannot be blank.", nil, false)
			}
		} else {
			a.Log("The username for the UCS System cannot be blank.", nil, false)
		}
	} else {
		a.Log("The URL for the UCS System cannot be blank.", nil, false)
	}
	return false
}

func (a *Application) addUCSPM(ip, username, password string) bool {
	if ip != "" {
		if username != "" {
			if password != "" {
				a.Config.Set("ucspm.url", ip)
				a.Config.Set("ucspm.username", username)
				a.Config.Set("ucspm.password", a.EncryptPassword(password))
				return true
			} else {
				a.Log("The password for the UCS Performance Manager system cannot be blank.", nil, false)
			}
		} else {
			a.Log("The username for the UCS Performance Manager system cannot be blank.", nil, false)
		}
	} else {
		a.Log("The URL for the UCS Performance Manager system cannot be blank.", nil, false)
	}
	return false
}

func (a *Application) addUCSPMSystem(ip, username, password string) {
	if !a.checkUCSPMExists(ip) {
		if a.addUCSPM(ip, username, password) {
			a.saveConfig()
			a.LogInfo("UCS Performance Manager system has been added successfully.", map[string]interface{}{"IP": ip, "Username": username}, false)
		}
	} else {
		a.LogInfo("A UCS Performance Manager already exists in the config file.", map[string]interface{}{"IP": ip, "Username": username}, false)
	}
}

func (a *Application) addUCSSystem(ip, username, password string) {
	if !a.checkUCSExists(ip) {
		if a.addUCS(ip, username, password) {
			a.saveConfig()
			a.LogInfo("New UCS system has been added successfully.", map[string]interface{}{"IP": ip, "Username": username}, false)
		} else {
			a.LogInfo("UCS System could not be added.", map[string]interface{}{"IP": ip, "Username": username}, false)
		}
	} else {
		a.LogInfo("A UCS System already exsists in the config file.", map[string]interface{}{"IP": ip, "Username": username}, false)
	}
}

func (a *Application) checkUCSExists(ip string) bool {
	a.Log("Search for UCS System in config file", map[string]interface{}{"IP": ip}, true)
	if a.Config.IsSet("ucs.systems") {
		a.getAllSystems()
		for i := 0; i < len(a.UCSSystems); i++ {
			if strings.TrimSpace(a.UCSSystems[i].ip) == strings.TrimSpace(ip) {
				return true
			}
		}
		return false
	}
	return false
}

func (a *Application) checkUCSPMExists(ip string) bool {
	a.getAllSystems()
	if a.Config.IsSet("ucspm.url") {
		return true
	}
	return false
}

func (a *Application) deleteUCS(ip string) bool {
	for i := 0; i < len(a.UCSSystems); i++ {
		if a.UCSSystems[i].ip == as.ToString(ip) {
			a.UCSSystems = append(a.UCSSystems[:i], a.UCSSystems[i+1:]...)
		}
	}
	return true
}

func (a *Application) deleteUCSSystem(ip string) {
	if a.checkUCSExists(ip) {
		if a.deleteUCS(ip) {
			a.saveConfig()
			a.LogInfo("UCS system has been deleted successfully.", map[string]interface{}{"IP": ip}, true)
		} else {
			a.Log("UCS System could not be deleted.", map[string]interface{}{"IP": ip}, false)
		}
	} else {
		a.LogInfo("UCS System does not exsists and so cannot be deleted.", map[string]interface{}{"IP": ip}, false)
	}
}

func (a *Application) deleteUCSPMSystem() {
	a.getAllSystems()
	if a.Config.IsSet("ucspm.url") {
		a.Config.Set("ucspm", false)
		a.saveConfig()
		a.LogInfo("UCS Performance Manager system has been deleted successfully.", nil, true)
	} else {
		a.LogInfo("UCS Performance Manager system does not exsists and so cannot be deleted.", nil, false)
	}
}

func (a *Application) getAllSystems() {
	tmp := as.ToSlice(a.Config.Get("ucs.systems"))
	a.Log("Located UCS Systems in the config file", map[string]interface{}{"Systems": len(tmp)}, true)
	a.readSystems(tmp)
}

func (a *Application) getAllUCSSystemsCount() int {
	tmp := as.ToSlice(a.Config.Get("ucs.systems"))
	return len(tmp)
}

func (a *Application) getAllUCSPMSystemsCount() int {
	if a.Config.IsSet("ucspm.url") {
		return 1
	} else {
		return 0
	}
}

func (a *Application) getEULAStatus() bool {
	if a.Config.IsSet("eula.agreed") {
		return true
	}
	return false
}

func (a *Application) processResponse(response string) {
	a.Log("Processing command line options.", map[string]interface{}{"args": response}, true)
	splits := strings.Split(response, "|")
	switch splits[0] {
	case "RUN":
		a.runAll(splits[1], splits[2])
	case "ADDUCS":
		a.addUCSSystem(splits[1], splits[2], splits[3])
	case "UPDATEUCS":
		a.updateUCSSystem(splits[1], splits[2], splits[3])
	case "DELETEUCS":
		a.deleteUCSSystem(splits[1])
	case "SHOWUCS":
		a.showUCSSystem(splits[1])
	case "SHOWALL":
		a.showUCSSystems()
	case "ADDUCSPM":
		a.addUCSPMSystem(splits[1], splits[2], splits[3])
	case "UPDATEUCSPM":
		a.updateUCSPMSystem(splits[1], splits[2], splits[3])
	case "DELETEUCSPM":
		a.deleteUCSPMSystem()
	case "SHOWUCSPM":
		a.showUCSPMSystem()
	case "SETINPUT":
		a.setInputFileName(splits[1])
	case "SETOUTPUT":
		a.setOutputFileName(splits[1])
	}
}
func (a *Application) readSystems(ucss []interface{}) bool {
	a.UCSSystems = nil
	for i := 0; i < len(ucss); i++ {
		var newlist map[string]string
		newlist = as.ToStringMapString(ucss[i])
		tmp := UCSSystemInfo{}
		tmp.ip = newlist["url"]
		tmp.username = newlist["username"]
		tmp.password = newlist["password"]
		a.UCSSystems = append(a.UCSSystems, tmp)
	}

	return true
}
func (a *Application) runAll(month, year string) {
	a.Log("Running inventory processes.", map[string]interface{}{"Month": month, "Year": year}, true)
	month, year = a.getReportDates(month, year)
	a.Log("Processed report dates.", map[string]interface{}{"Month": month, "Year": year}, true)
	a.Report.Month = month
	a.Report.Year = year
	a.RunStage2()
}

func (a *Application) setInputFileName(filename string) {
	a.Config.Set("input.file", filename)
	a.saveConfig()
}

func (a *Application) setOutputFileName(filename string) {
	a.Config.Set("output.file", filename)
	a.saveConfig()
}

func (a *Application) showUCS(ip string) {
	for i := 0; i < len(a.UCSSystems); i++ {
		if a.UCSSystems[i].ip == as.ToString(ip) {
			a.LogInfo("UCS Domain", map[string]interface{}{"URL": a.UCSSystems[i].ip}, false)
			a.LogInfo("UCS Domain", map[string]interface{}{"Username": a.UCSSystems[i].username}, false)
			a.LogInfo("UCS Domain", map[string]interface{}{"Password": a.UCSSystems[i].password}, false)
		}
	}
}

func (a *Application) showUCSSystem(ip string) {
	if a.checkUCSExists(ip) {
		a.showUCS(ip)
	} else {
		a.Log("The UCS Domain does not exist and so cannot be displayed.", map[string]interface{}{"URL": ip}, false)
	}
}

func (a *Application) showUCSSystems() {
	a.getAllSystems()
	for i := 0; i < len(a.UCSSystems); i++ {
		a.LogInfo("UCS Domain", map[string]interface{}{"URL": a.UCSSystems[i].ip}, false)
		a.LogInfo("UCS Domain", map[string]interface{}{"Username": a.UCSSystems[i].username}, false)
		a.LogInfo("UCS Domain", map[string]interface{}{"Password": a.UCSSystems[i].password}, false)
	}
	a.showUCSPMSystem()
}

func (a *Application) showUCSPMSystem() {
	if a.Config.IsSet("ucspm.url") {
		a.LogInfo("UCS Performance Manager", map[string]interface{}{"URL": a.Config.GetString("ucspm.url")}, false)
		a.LogInfo("UCS Performance Manager", map[string]interface{}{"Username": a.Config.GetString("ucspm.username")}, false)
		a.LogInfo("UCS Performance Manager", map[string]interface{}{"Password": a.Config.GetString("ucspm.password")}, false)
	}
}

func (a *Application) updateUCS(ip, username, password string) bool {
	for i := 0; i < len(a.UCSSystems); i++ {
		if a.UCSSystems[i].ip == as.ToString(ip) {
			a.UCSSystems[i].username = username
			a.UCSSystems[i].password = a.EncryptPassword(password)
		}
	}
	return true
}

func (a *Application) updateUCSSystem(ip, username, password string) {
	if a.checkUCSExists(ip) {
		if a.updateUCS(ip, username, password) {
			a.saveConfig()
			a.LogInfo("Update to UCS system has been completed successfully.", map[string]interface{}{"IP": ip, "Username": username}, false)
		} else {
			a.LogInfo("UCS System could not be updated.", map[string]interface{}{"IP": ip, "Username": username}, false)
		}
	} else {
		a.LogInfo("UCS System does not exsist and can therefore not be updated.", map[string]interface{}{"IP": ip, "Username": username}, false)
	}
}

func (a *Application) updateUCSPMSystem(ip, username, password string) {
	if a.checkUCSPMExists(ip) {
		if a.addUCSPM(ip, username, password) {
			a.saveConfig()
			a.LogInfo("UCS Performance Manager system has been updated successfully.", map[string]interface{}{"IP": ip, "Username": username}, false)
		}
	} else {
		a.LogInfo("A UCS Performance Manager instance does not exist in the config file.", map[string]interface{}{"IP": ip, "Username": username}, false)
	}
}