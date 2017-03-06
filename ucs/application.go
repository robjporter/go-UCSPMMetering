package ucs

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"../flags"

	functions "github.com/robjporter/go-functions"
	"github.com/robjporter/go-functions/as"
	"github.com/robjporter/go-functions/banner"
	"github.com/robjporter/go-functions/colors"
	"github.com/robjporter/go-functions/http"
	"github.com/robjporter/go-functions/terminal"
	"github.com/spf13/viper"
	yaml "pkg.re/yaml.v2"
)

var steps = []string{"downloading source", "installing deps", "compiling", "packaging", "seeding database", "deploying", "staring servers"}

type system struct {
	ip       string
	username string
	password string
	cookie   string
	name     string
	version  string
}

type match struct {
	serverposition string
	serverserial   string
	serveruuid     string
	servername     string
	serverpid      string
	serverdn       string
	serverdescr    string
	servermodel    string
	serverouuid    string
	ucsname        string
	ucsversion     string
	ucsip          string
}

type Application struct {
	configFile string
	uuid       []string
	systems    []system
	matches    []match
	matched    []match
	unmatched  []string
}

// ===========PUBLIC============

func (a *Application) CheckConfig(filename string) bool {
	configName := ""
	configExtension := ""
	configPath := ""

	splits := strings.Split(filepath.Base(filename), ".")
	if len(splits) == 2 {
		configName = splits[0]
		configExtension = splits[1]
	}
	configPath = filepath.Dir(filename)

	file := filename
	if functions.Exists(file) {
		viper.SetConfigName(configName)
		viper.SetConfigType(configExtension)
		viper.AddConfigPath(configPath)
		a.configFile = file
		return true
	} else {
		return false
	}
}

func (a *Application) Run() {
	a.loadConfig()
	a.start()
}

// ===========PRIVATE============

func (a *Application) start() {
	a.drawHeader()
	if viper.IsSet("input.file") {
		if a.readUUIDFile(viper.GetString("input.file")) {
			if viper.IsSet("ucs.systems") {
				tmp := as.ToSlice(viper.Get("ucs.systems"))
				a.readSystems(tmp)
				a.displayInfo()
				a.processInfo()
			}
		} else {
			fmt.Println("Input file was not found.  Please check and try again.")
		}
	} else {
		fmt.Println("CONFIG FILE MISSING INPUT FILE")
		os.Exit(0)
	}
}

func (a *Application) loadConfig() {
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}

func (a *Application) saveConfig() {
	out, err := yaml.Marshal(viper.AllSettings())
	if err == nil {
		fp, err := os.Create(a.configFile)
		if err == nil {
			defer fp.Close()
			_, err = fp.Write(out)
		}
	}
}

func (a *Application) drawHeader() {
	terminal.ClearScreen()
	banner.PrintNewFigure("UCS", "rounded", true)
	fmt.Println(colors.Color("Cisco Unified Computing System", colors.BRIGHTYELLOW))
	banner.BannerPrintLineS("=", 60)
	fmt.Println("\n\n")
}

func (a *Application) readUUIDFile(filename string) bool {
	if functions.Exists(filename) {

		var uuids = viper.New()
		uuids.SetConfigType("json")
		uuids.SetConfigName(functions.GetFilenameNoExtension(filename))
		uuids.AddConfigPath("./")
		err := uuids.ReadInConfig()

		if err == nil {
			if uuids.IsSet("uuids") {
				a.uuid = uuids.GetStringSlice("uuids")
			}
		}

		return true
	}
	return false
}

func (a *Application) displayInfo() {
	fmt.Printf("Found %d UUID's and %d UCS systems to scan and match.\n", len(a.uuid), len(a.systems))
}

func (a *Application) readSystems(systems []interface{}) bool {
	var newlist map[string]string
	for i := 0; i < len(systems); i++ {
		newlist = as.ToStringMapString(systems[i])
		tmp := system{}
		tmp.ip = newlist["url"]
		tmp.username = newlist["username"]
		tmp.password = flags.DecryptPassword(newlist["password"])
		a.systems = append(a.systems, tmp)
	}
	return true
}

func (a *Application) makeConnectionUrl(position int) {
	if !strings.Contains(a.systems[position].ip, "http") {
		a.systems[position].ip = "https://" + a.systems[position].ip
	}
	if !strings.Contains(a.systems[position].ip, "nuova") {
		a.systems[position].ip = a.systems[position].ip + "/nuova"
	}
}

func (a *Application) connectToUCS(sys system) (string, string) {
	xml, err := getLoginXML()
	headers := make(map[string]string)
	result := ""
	version := ""
	if err == nil {
		headers["Content-Type"] = "application/xml"
		xml = replaceString(xml, "|USERNAME|", sys.username)
		xml = replaceString(xml, "|PASSWORD|", sys.password)
		code, response, err := http.SendUnsecureHTTPSRequest(sys.ip, "POST", xml, headers)
		if err == nil {
			if code == 200 {
				result = GetQueryResponseData(response, "aaaLogin", "", "outCookie")
				version = GetQueryResponseData(response, "aaaLogin", "", "outVersion")
			}
		} else {
			fmt.Println(err)
		}
	}
	return result, version
}

func (a *Application) getUCSData(sys system) {
	fmt.Println("Getting: ", sys.ip)
	result := a.getServerInfo(sys)

	getBladeDNs := getServerDN(result)
	for i := 0; i < len(getBladeDNs); i++ {
		result := a.getUCSServerDetail(getBladeDNs[i], sys)
		var mat match
		model, name, ouuid, pid, serial, uuid, position, description := getServerDetail(result)
		mat.serverdn = getBladeDNs[i]
		mat.servermodel = model
		mat.serverdescr = description
		mat.servername = name
		mat.serverpid = pid
		mat.serverposition = a.formatUCSPosition(position)
		mat.serverserial = serial
		mat.serveruuid = uuid
		mat.serverouuid = ouuid
		mat.ucsname = sys.name
		mat.ucsversion = sys.version
		a.matches = append(a.matches, mat)
	}
}

func (a *Application) formatUCSPosition(pos string) string {
	result := ""
	if pos != "" {
		splits := strings.Split(pos, "/")
		if len(splits) == 2 {
			result = "Chassis: " + splits[0] + " | Blade: " + splits[1]
		} else {
			result = pos
		}
	}
	return result
}

func (a *Application) getUCSServerDetail(dn string, sys system) string {
	xml, err := getServerDetailXML()
	headers := make(map[string]string)
	result := ""
	if err == nil {
		headers["Content-Type"] = "application/xml"
		xml = replaceString(xml, "|COOKIE|", sys.cookie)
		xml = replaceString(xml, "|DN|", dn)
		code, response, err := http.SendUnsecureHTTPSRequest(sys.ip, "POST", xml, headers)
		if err == nil {
			if code == 200 {
				result = response
			}
		} else {
			fmt.Println(err)
		}
	}
	return result
}

func (a *Application) getServerInfo(sys system) string {
	xml, err := GetAllServers()
	result := ""
	headers := make(map[string]string)

	if err == nil {
		headers["Content-Type"] = "application/xml"
		xml = replaceString(xml, "|COOKIE|", sys.cookie)

		code, response, err := http.SendUnsecureHTTPSRequest(sys.ip, "POST", xml, headers)
		if err == nil {
			if code == 200 {
				result = response
			}
		} else {
			fmt.Println(err)
		}
	}

	return result
}

func (a *Application) makeUCSConnections() {
	for i := 0; i < len(a.systems); i++ {
		cookie, version := a.connectToUCS(a.systems[i])
		if cookie != "" {
			a.systems[i].cookie = cookie
			a.systems[i].version = version
			a.systems[i].name = a.getUCSName(a.systems[i])
		}
	}
}

func (a *Application) getAllUUIDInfo() {
	for i := 0; i < len(a.systems); i++ {
		a.getUCSData(a.systems[i])
	}
}

func (a *Application) getUCSName(sys system) string {
	xml, err := getSystemDetailXML()
	headers := make(map[string]string)
	result := ""
	if err == nil {
		headers["Content-Type"] = "application/xml"
		xml = replaceString(xml, "|COOKIE|", sys.cookie)
		code, response, err := http.SendUnsecureHTTPSRequest(sys.ip, "POST", xml, headers)
		if err == nil {
			if code == 200 {
				result = GetQueryResponseData2(response, "configResolveClass", "outConfigs", "topSystem", "name")
			}
		} else {
			fmt.Println(err)
		}
	}
	return result
}

func (a *Application) processMatchedUUID() {
	fmt.Printf("Matching %d UUID(s) with %d server(s) from %d domain(s)\n", len(a.uuid), len(a.matches), len(a.systems))
	unmatched := make([]string, len(a.uuid))
	copy(unmatched, a.uuid)
	for i := 0; i < len(a.matches); i++ {
		for j := 0; j < len(a.uuid); j++ {
			if a.matches[i].serveruuid == a.uuid[j] || a.matches[i].serverouuid == a.uuid[j] {
				a.matched = append(a.matched, a.matches[i])
				unmatched[j] = "REMOVE"
			}
		}
	}
	a.removeMatched(unmatched)
	fmt.Printf("Successfully matched %d of %d UUID.\n", len(a.matched), len(a.uuid))
	if len(a.matched) < len(a.uuid) {
		fmt.Println("The unmatched UUID are: ", a.unmatched)
	}
}

func (a *Application) removeMatched(list []string) {
	for i := 0; i < len(list); i++ {
		if strings.TrimSpace(list[i]) != "REMOVE" {
			a.unmatched = append(a.unmatched, list[i])
		}
	}
}

func (a *Application) logoutDomains() {
	for i := 0; i < len(a.systems); i++ {
		a.logoutDomain(a.systems[i])
	}
}

func (a *Application) logoutDomain(sys system) bool {
	xml, err := getLogoutXML()
	headers := make(map[string]string)
	result := false
	if err == nil {
		headers["Content-Type"] = "application/xml"
		xml = replaceString(xml, "|COOKIE|", sys.cookie)
		code, _, err := http.SendUnsecureHTTPSRequest(sys.ip, "POST", xml, headers)
		if err == nil {
			if code == 200 {
				result = true
			}
		} else {
			fmt.Println(err)
		}
	}
	return result
}

func (a *Application) exportToCSV() {
	csv := "name,uuid,serial,domain,domainversion,position,model,pid,description,dn\n"
	for i := 0; i < len(a.matched); i++ {
		csv += a.matched[i].servername + "," + a.matched[i].serveruuid + "," + a.matched[i].serverserial + ","
		csv += a.matched[i].ucsname + "," + a.matched[i].ucsversion + "," + a.matched[i].serverposition + "," + a.matched[i].servermodel + ","
		csv += a.matched[i].serverpid + "," + a.matched[i].serverdescr + "," + a.matched[i].serverdn + "\n"

	}
	if viper.IsSet("output.file") {
		f, err := os.Create(viper.GetString("output.file"))
		if err == nil {
			_, err := f.Write([]byte(csv))
			if err == nil {
				fmt.Println("\nFile saved successfully.")
			} else {
				fmt.Println("\nThere was an error saving the file: ", err)
			}
		}
		defer f.Close()
	}
}

func (a *Application) processInfo() {
	for i := 0; i < len(a.systems); i++ {
		a.makeConnectionUrl(i)
	}
	a.makeUCSConnections()
	a.getAllUUIDInfo()
	a.logoutDomains()
	a.processMatchedUUID()
	a.exportToCSV()
}
