package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func main() {
	if len(os.Args) == 6 {
		ip := os.Args[1]
		username := os.Args[2]
		password := os.Args[3]
		vCenter := os.Args[4]
		serverNumber := os.Args[5]

		json := []byte(`{"end": 1490266161,"series": true,"start": 1485907200,"downsample": "1h-avg","metrics": [{"metric": "` + vCenter + `/cpuUsage_cpuUsage","rate": false,"emit": false,"rateOptions": {},"aggregator": "avg","tags": {"key": ["Devices/` + vCenter + `/datacenters/Datacenter_datacenter-21/hosts/HostSystem_host-` + serverNumber + `"]},"name": "Usage-raw"},{"name": "Usage","expression": "rpn:Usage-raw,100,/"}],"returnset": "EXACT","tags": {}}`)

		body := bytes.NewBuffer(json)

		// Create client
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client := &http.Client{Transport: tr}

		// Create request
		req, err := http.NewRequest("POST", "https://"+ip+"/api/performance/query/", body)
		req.SetBasicAuth(username, password)

		// Headers
		req.Header.Add("Accept-Charset", "utf-8")
		req.Header.Add("Content-type", "application/json")

		// Fetch Request
		resp, err := client.Do(req)

		if err != nil {
			fmt.Println("Failure : ", err)
		}

		// Read Response Body
		respBody, _ := ioutil.ReadAll(resp.Body)

		// Display Results
		fmt.Println("response Status : ", resp.Status, "\n")
		fmt.Println("response Headers : ", resp.Header, "\n")
		fmt.Println("response Body : ", string(respBody))

		contents := "response Status : " + string(resp.Status) + "\n\n"
		for k, v := range resp.Header {
			contents += "response Header Key:" + k + " | Value: " + v[0] + "\n"
		}
		contents += "\n"
		contents += "response Body : " + string(respBody)

		t := time.Now()
		filename := "HostSystem_host-" + serverNumber + "-" + t.Format("20060102150405") + ".txt"

		ioutil.WriteFile(filename, []byte(contents), 0644)
	}
}
