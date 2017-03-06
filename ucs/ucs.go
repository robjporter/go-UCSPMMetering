package ucs

import (
	"fmt"
	"strings"

	"github.com/beevik/etree"
)

func getLoginXML() (string, error) {
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

func GetQueryResponseData2(xml string, root string, element string, subelement string, attribute string) string {
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

func GetQueryResponseData(xml string, root string, element string, attribute string) string {
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

func GetAllServers() (string, error) {
	doc := etree.NewDocument()
	doc.CreateProcInst("xml", `version="1.0" encoding="UTF-8"`)
	doc.CreateProcInst("xml-stylesheet", `type="text/xsl" href="style.xsl"`)

	login := doc.CreateElement("configFindDnsByClassId")
	login.CreateAttr("classId", "computeItem")
	login.CreateAttr("cookie", "|COOKIE|")
	doc.Indent(2)
	return doc.WriteToString()
}

func GetChassisXML() (string, error) {
	doc := etree.NewDocument()
	doc.CreateProcInst("xml", `version="1.0" encoding="UTF-8"`)
	doc.CreateProcInst("xml-stylesheet", `type="text/xsl" href="style.xsl"`)

	resolve := doc.CreateElement("configResolveDns")
	resolve.CreateAttr("cookie", "|COOKIE|")
	resolve.CreateAttr("inHierarchical", "false")
	dns := resolve.CreateElement("inDns")
	dn := dns.CreateElement("dn")
	dn.CreateAttr("value", "sys/chassis-1")
	doc.Indent(2)
	return doc.WriteToString()
}

func getServerDN(result string) []string {
	doc := etree.NewDocument()
	results := []string{}
	if err := doc.ReadFromString(result); err != nil {
		panic(err)
	}
	root := doc.SelectElement("configFindDnsByClassId")
	out := root.SelectElement("outDns")

	for _, server := range out.SelectElements("dn") {
		dn := server.SelectAttrValue("value", "unknown")
		results = append(results, dn)
	}
	return results
}

func getServerDetailXML() (string, error) {
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

func getSystemDetailXML() (string, error) {
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

func getServerDetail(result string) (string, string, string, string, string, string, string, string) {
	doc := etree.NewDocument()
	if err := doc.ReadFromString(result); err != nil {
		panic(err)
	}
	root := doc.SelectElement("configResolveDn")
	out := root.SelectElement("outConfig")

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
