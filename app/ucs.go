package app

import (
	"fmt"
	"strings"

	"github.com/robjporter/go-functions/etree"
	"github.com/robjporter/go-functions/http"
)

func (a *Application) ucsExportToCSV() {
	csv := "name,uuid,serial,domain,domainversion,position,model,pid,description,dn\n"
	for i := 0; i < len(a.UCS.Matched); i++ {
		csv += a.UCS.Matched[i].servername + "," + a.UCS.Matched[i].serveruuid + "," + a.UCS.Matched[i].serverserial + ","
		csv += a.UCS.Matched[i].ucsname + "," + a.UCS.Matched[i].ucsversion + "," + a.UCS.Matched[i].serverposition + "," + a.UCS.Matched[i].servermodel + ","
		csv += a.UCS.Matched[i].serverpid + "," + a.UCS.Matched[i].serverdescr + "," + a.UCS.Matched[i].serverdn + "\n"

	}
	if a.Config.IsSet("output.file") {
		a.saveFile(a.Config.GetString("output.file"), csv)
	}
}

func (a *Application) ucsProcessMatchedUUID() {
	a.LogInfo("Starting UUID Match process.", map[string]interface{}{"UUID": len(a.UCS.UUID), "Servers": len(a.UCS.Matches), "Domains": len(a.UCS.Systems)}, true)
	unmatched := make([]string, len(a.UCS.UUID))
	copy(unmatched, a.UCS.UUID)
	for i := 0; i < len(a.UCS.Matches); i++ {
		matched := false
		for j := 0; j < len(a.UCS.UUID); j++ {
			if !matched {
				if a.UCS.Matches[i].serveruuid == a.UCS.UUID[j] {
					a.UCS.Matched = append(a.UCS.Matched, a.UCS.Matches[i])
					unmatched[j] = "REMOVE"
					matched = true
				} else if a.UCS.Matches[i].serverouuid == a.UCS.UUID[j] {
					a.UCS.Matches[i].serveruuid = a.UCS.Matches[i].serverouuid
					a.UCS.Matched = append(a.UCS.Matched, a.UCS.Matches[i])
					unmatched[j] = "REMOVE"
					matched = true
				} else if a.transposeUUID(a.UCS.Matches[i].serveruuid) == a.UCS.UUID[j] {
					a.UCS.Matches[i].serveruuid = a.transposeUUID(a.UCS.Matches[i].serveruuid)
					a.UCS.Matched = append(a.UCS.Matched, a.UCS.Matches[i])
					unmatched[j] = "REMOVE"
					matched = true
				}
			}
		}
	}
	a.ucsRemoveMatched(unmatched)
}

func (a *Application) transposeUUID(uuid string) string {
	splits := strings.Split(uuid, "-")
	if len(splits) == 5 {
		splits[0] = rotateNumber(splits[0])
		splits[1] = rotateNumber(splits[1])
		splits[2] = rotateNumber(splits[2])
		uuid = splits[0] + "-" + splits[1] + "-" + splits[2] + "-" + splits[3] + "-" + splits[4]
	}
	return uuid
}

func rotateNumber(number string) string {
	if len(number) == 4 {
		part1 := number[0:2]
		part2 := number[2:4]
		number = part2 + part1
	} else if len(number) == 8 {
		part1 := number[0:2]
		part2 := number[2:4]
		part3 := number[4:6]
		part4 := number[6:8]
		number = part4 + part3 + part2 + part1
	}
	return number
}

func (a *Application) ucsRemoveMatched(list []string) {
	for i := 0; i < len(list); i++ {
		if strings.TrimSpace(list[i]) != "REMOVE" {
			a.UCS.Unmatched = append(a.UCS.Unmatched, list[i])
		}
	}
}

func (a *Application) ucsReadUCSPMUUIDFile() {
	/*filename := "data/Stage4-DiscoveredUUID.json"
	if functions.Exists(filename) {

		var uuids = viper.New()
		uuids.SetConfigType("json")
		uuids.SetConfigName(functions.GetFilenameNoExtension(filename))
		uuids.AddConfigPath("./")
		err := uuids.ReadInConfig()

		if err == nil {
			if uuids.IsSet("uuids") {
				a.UCS.UUID = uuids.GetStringSlice("uuids")
			}
		}
		a.LogInfo("Loaded UCS UUID from UCS Performance Manager.", map[string]interface{}{"Count": len(a.UCS.UUID), "UUID": a.UCS.UUID}, true)
	}*/
	a.UCS.UUID = a.UCSPM.ProcessedUUID
}

func (a *Application) ucsInit() {

}

func (a *Application) ucsInventory() {
	a.LogInfo("Preparing to run inventory on UCS Managers.", map[string]interface{}{"UniqueUUID": len(a.UCSPM.ProcessedUUID), "UCSSystems": len(a.UCS.Systems)}, false)
	for i := 0; i < len(a.UCS.Systems); i++ {
		a.ucsMakeConnectionURL(i)
	}
	a.ucsConnection()
	a.ucsGetAllUUIDInfo()
	a.ucsLogoutDomains()
	a.ucsReadUCSPMUUIDFile()
	a.ucsProcessMatchedUUID()
	a.ucsExportToCSV()
}

func (a *Application) ucsConnection() {
	for i := 0; i < len(a.UCS.Systems); i++ {
		cookie, version := a.ucsConnectToSystem(a.UCS.Systems[i])
		if cookie != "" {
			a.UCS.Systems[i].cookie = cookie
			a.UCS.Systems[i].version = version
			a.UCS.Systems[i].name = a.ucsGetSystemName(a.UCS.Systems[i])
			a.LogInfo("Successfully connected to UCS System.", map[string]interface{}{"URL": a.UCS.Systems[i].ip, "Name": a.UCS.Systems[i].name, "Version": a.UCS.Systems[i].version}, false)
		} else {
			a.Log("Failed to connect to UCS System.", map[string]interface{}{"URL": a.UCS.Systems[i].ip}, true)
		}
	}
}

func (a *Application) ucsConnectToSystem(sys UCSSystemInfo) (string, string) {
	xml, err := ucsGetLoginXML()
	headers := make(map[string]string)
	result := ""
	version := ""
	if err == nil {
		headers["Content-Type"] = "application/xml"
		xml = replaceString(xml, "|USERNAME|", sys.username)
		xml = replaceString(xml, "|PASSWORD|", a.DecryptPassword(sys.password))
		code, response, err := http.SendUnsecureHTTPSRequest(sys.ip, "POST", xml, headers)

		a.addCommand(sys.ip, xml, headers, response, code, err)

		if err == nil {
			if code == 200 {
				result = getQueryResponseData(response, "aaaLogin", "", "outCookie")
				version = getQueryResponseData(response, "aaaLogin", "", "outVersion")
			} else {
				a.Log("Failed to gain a 200 response code.", map[string]interface{}{"URL": sys.ip, "Code": code, "Response": response}, true)
			}
		} else {
			a.Log("Failed to connect to UCS System.", map[string]interface{}{"URL": sys.ip, "ERROR": err}, true)
		}
	}
	return result, version
}

func (a *Application) ucsGetSystemName(sys UCSSystemInfo) string {
	xml, err := ucsGetSystemDetailXML()
	headers := make(map[string]string)
	result := ""
	if err == nil {
		headers["Content-Type"] = "application/xml"
		xml = replaceString(xml, "|COOKIE|", sys.cookie)
		code, response, err := http.SendUnsecureHTTPSRequest(sys.ip, "POST", xml, headers)

		a.addCommand(sys.ip, xml, headers, response, code, err)

		if err == nil {
			if code == 200 {
				result = getQueryResponseData2(response, "configResolveClass", "outConfigs", "topSystem", "name")
			} else {
				a.Log("Failed to gain a 200 response code.", map[string]interface{}{"URL": sys.ip, "Code": code, "Response": response}, true)
			}
		} else {
			a.Log("Failed to connect to UCS System.", map[string]interface{}{"URL": sys.ip, "ERROR": err}, true)
		}
	}
	return result
}

func (a *Application) ucsLogoutDomain(sys UCSSystemInfo) bool {
	xml, err := getLogoutXML()
	headers := make(map[string]string)
	result := false
	if err == nil {
		headers["Content-Type"] = "application/xml"
		xml = replaceString(xml, "|COOKIE|", sys.cookie)
		code, response, err := http.SendUnsecureHTTPSRequest(sys.ip, "POST", xml, headers)

		a.addCommand(sys.ip, xml, headers, response, code, err)

		if err == nil {
			if code == 200 {
				result = true
				a.Log("Succesfully logged out of UCS Domain.", map[string]interface{}{"URL": sys.ip, "Code": code, "Response": response}, true)
			} else {
				a.Log("Failed to gain a 200 response code.", map[string]interface{}{"URL": sys.ip, "Code": code, "Response": response}, true)
			}
		} else {
			a.Log("Failed to connect to UCS System.", map[string]interface{}{"URL": sys.ip, "ERROR": err}, true)
		}
	}
	return result
}

func (a *Application) ucsLogoutDomains() {
	a.LogInfo("Logging out of all UCS Domains.", nil, true)
	for i := 0; i < len(a.UCS.Systems); i++ {
		a.ucsLogoutDomain(a.UCS.Systems[i])
	}
}

func (a *Application) ucsMakeConnectionURL(position int) {
	tmpURL := a.UCS.Systems[position].ip
	if !strings.Contains(a.UCS.Systems[position].ip, "http") {
		a.UCS.Systems[position].ip = "https://" + a.UCS.Systems[position].ip
	}
	if !strings.Contains(a.UCS.Systems[position].ip, "nuova") {
		a.UCS.Systems[position].ip = a.UCS.Systems[position].ip + "/nuova"
	}
	a.Log("Changing UCS System connection URL.", map[string]interface{}{"Original": tmpURL, "Corrected": a.UCS.Systems[position].ip}, true)
}

func (a *Application) ucsGetAllUUIDInfo() {
	a.LogInfo("Getting all UCS System UUID Inventory.", nil, true)
	for i := 0; i < len(a.UCS.Systems); i++ {
		a.ucsGetUCSData(a.UCS.Systems[i])
	}
}

func (a *Application) ucsGetUCSData(sys UCSSystemInfo) {
	result := a.ucsGetServerInfo(sys)
	getBladeDNs := ucsGetServerDN(result)
	a.LogInfo("Getting UCS System Server Detail.", map[string]interface{}{"URL": sys.ip, "Servers": len(getBladeDNs)}, true)
	for i := 0; i < len(getBladeDNs); i++ {
		result := a.ucsGetServerDetail(getBladeDNs[i], sys)
		var mat UCSSystemMatchInfo
		model, name, ouuid, pid, serial, uuid, position, description := getServerDetail(result)
		mat.serverdn = getBladeDNs[i]
		mat.servermodel = model
		mat.serverdescr = description
		mat.servername = name
		mat.serverpid = pid
		mat.serverposition = ucsFormatPosition(position)
		mat.serverserial = serial
		mat.serveruuid = uuid
		mat.serverouuid = ouuid
		mat.ucsname = sys.name
		mat.ucsversion = sys.version
		a.UCS.Matches = append(a.UCS.Matches, mat)
	}
}

func (a *Application) ucsGetUCSSystem(uuid string) UCSSystemMatchInfo {
	for i := 0; i < len(a.UCS.Matched); i++ {
		if a.UCS.Matched[i].serveruuid == uuid {
			return a.UCS.Matched[i]
		}
	}
	return UCSSystemMatchInfo{}
}

/*
	UCS HELPER FUNCTIONS
*/

func (a *Application) ucsGetServerDetail(dn string, sys UCSSystemInfo) string {
	xml, err := ucsGetServerDetailXML()
	headers := make(map[string]string)
	result := ""
	if err == nil {
		headers["Content-Type"] = "application/xml"
		xml = replaceString(xml, "|COOKIE|", sys.cookie)
		xml = replaceString(xml, "|DN|", dn)
		code, response, err := http.SendUnsecureHTTPSRequest(sys.ip, "POST", xml, headers)

		a.addCommand(sys.ip, xml, headers, response, code, err)

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

func (a *Application) ucsGetServerInfo(sys UCSSystemInfo) string {
	xml, err := ucsGetAllServersXML()
	result := ""
	headers := make(map[string]string)

	if err == nil {
		headers["Content-Type"] = "application/xml"
		xml = replaceString(xml, "|COOKIE|", sys.cookie)

		code, response, err := http.SendUnsecureHTTPSRequest(sys.ip, "POST", xml, headers)

		a.addCommand(sys.ip, xml, headers, response, code, err)

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

func ucsGetServerDN(result string) []string {
	doc := etree.NewDocument()
	results := []string{}
	if err := doc.ReadFromString(result); err != nil {
		panic(err)
	}
	root := doc.SelectElement("configFindDnsByClassId")
	if root != nil {
		out := root.SelectElement("outDns")
		if out != nil {
			for _, server := range out.SelectElements("dn") {
				dn := server.SelectAttrValue("value", "unknown")
				results = append(results, dn)
			}
		}
	}
	return results
}

func getServerDetail(result string) (string, string, string, string, string, string, string, string) {
	doc := etree.NewDocument()
	if err := doc.ReadFromString(result); err != nil {
		panic(err)
	}
	root := doc.SelectElement("configResolveDn")
	if root != nil {
		out := root.SelectElement("outConfig")
		if out != nil {
			var server *etree.Element
			if server = out.SelectElement("computeRackUnit"); server == nil {
				if server = out.SelectElement("computeBlade"); server == nil {
					fmt.Println("ERROR")
				}
			}

			model := server.SelectAttrValue("model", "unknown")
			name := server.SelectAttrValue("name", "unknown")
			ouuid := server.SelectAttrValue("originalUuid", "unknown")
			pid := server.SelectAttrValue("partNumber", "unknown")
			serial := server.SelectAttrValue("serial", "unknown")
			uuid := server.SelectAttrValue("uuid", "unknown")
			position := server.SelectAttrValue("serverId", "unknown")
			decription := server.SelectAttrValue("descr", "unknown")
			return model, name, ouuid, pid, serial, uuid, position, decription
		}
	}
	return "", "", "", "", "", "", "", ""
}

func ucsFormatPosition(pos string) string {
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

func getQueryResponseData(xml string, root string, element string, attribute string) string {
	results := ""
	doc := etree.NewDocument()
	if err := doc.ReadFromString(xml); err == nil {
		root := doc.SelectElement(root)
		if element == "" {
			return root.SelectAttrValue(attribute, "unknown")
		} else {
			element := root.SelectElement(element)
			if attribute == "" {
			} else {
				return element.SelectAttrValue(attribute, "unknown")
			}
		}
	}
	return results
}

func getQueryResponseData2(xml string, root string, element string, subelement string, attribute string) string {
	results := ""
	doc := etree.NewDocument()
	if err := doc.ReadFromString(xml); err == nil {
		root2 := doc.SelectElement(root)
		if root2 != nil {
			out := root2.SelectElement(element)
			if out != nil {
				done := false
				for _, book := range out.SelectElements(subelement) {
					if !done {
						results = book.SelectAttrValue(attribute, "")
						done = true
					}
				}
			}
		}
	}
	return results
}

func ucsGetSystemDetailXML() (string, error) {
	doc := etree.NewDocument()
	doc.CreateProcInst("xml", `version="1.0" encoding="UTF-8"`)
	doc.CreateProcInst("xml-stylesheet", `type="text/xsl" href="style.xsl"`)

	resolve := doc.CreateElement("configResolveClass")
	resolve.CreateAttr("cookie", "|COOKIE|")
	resolve.CreateAttr("classId", "topSystem")
	resolve.CreateAttr("inHierarchical", "false")

	doc.Indent(2)
	return doc.WriteToString()
}

func ucsGetServerDetailXML() (string, error) {
	doc := etree.NewDocument()
	doc.CreateProcInst("xml", `version="1.0" encoding="UTF-8"`)
	doc.CreateProcInst("xml-stylesheet", `type="text/xsl" href="style.xsl"`)

	resolve := doc.CreateElement("configResolveDn")
	resolve.CreateAttr("dn", "|DN|")
	resolve.CreateAttr("cookie", "|COOKIE|")
	resolve.CreateAttr("inHierarchical", "false")
	doc.Indent(2)
	return doc.WriteToString()
}

func ucsGetAllServersXML() (string, error) {
	doc := etree.NewDocument()
	doc.CreateProcInst("xml", `version="1.0" encoding="UTF-8"`)
	doc.CreateProcInst("xml-stylesheet", `type="text/xsl" href="style.xsl"`)

	login := doc.CreateElement("configFindDnsByClassId")
	login.CreateAttr("classId", "computeItem")
	login.CreateAttr("cookie", "|COOKIE|")
	doc.Indent(2)
	return doc.WriteToString()
}

func ucsGetLoginXML() (string, error) {
	doc := etree.NewDocument()
	doc.CreateProcInst("xml", `version="1.0" encoding="UTF-8"`)
	doc.CreateProcInst("xml-stylesheet", `type="text/xsl" href="style.xsl"`)

	login := doc.CreateElement("aaaLogin")
	login.CreateAttr("inName", "|USERNAME|")
	login.CreateAttr("inPassword", "|PASSWORD|")
	doc.Indent(2)
	return doc.WriteToString()
}

func getLogoutXML() (string, error) {
	doc := etree.NewDocument()
	doc.CreateProcInst("xml", `version="1.0" encoding="UTF-8"`)
	doc.CreateProcInst("xml-stylesheet", `type="text/xsl" href="style.xsl"`)

	login := doc.CreateElement("aaaLogout")
	login.CreateAttr("inCookie", "|COOKIE|")
	doc.Indent(2)
	return doc.WriteToString()
}

func replaceString(xml string, param string, value string) string {
	return strings.Replace(xml, param, value, -1)
}
