package app

import (
	"github.com/robjporter/go-functions/logrus"
	"github.com/robjporter/go-functions/viper"
)

type UCSSystemInfo struct {
	ip       string
	username string
	password string
	cookie   string
	name     string
	version  string
}

type UCSSystemMatchInfo struct {
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

type CommandInfo struct {
	RequestURL     string
	RequestHeaders map[string]string
	RequestBody    string
	ResponseBody   string
	ResponseCode   int
	ResponseError  string
}

type UCSInfo struct {
	configFile string
	UUID       []string
	Systems    []UCSSystemInfo
	Matches    []UCSSystemMatchInfo
	Matched    []UCSSystemMatchInfo
	Unmatched  []string
}

type UCSPMInfo struct {
	Routers       map[string]string
	TidCount      int
	Devices       []UCSPMDeviceInfo
	host          string
	username      string
	password      string
	ProcessedUUID []string
}

type ReportInfo struct {
	Month string
	Year  string
}

type AppStatus struct {
	eula       bool
	ucsCount   int
	ucspmCount int
}

type CombinedResults struct {
	ucspmName   string
	ucspmUID    string
	ucspmKey    string
	ucspmUUID   string
	ucsName     string
	ucsPosition string
	ucsSerial   string
	ucsDN       string
	ucsDesc     string
	ucsModel    string
	ucsSystem   string
	isManaged   bool
	reportData  []ReportData
}

type Application struct {
	ConfigFile string
	Debug      bool
	Config     *viper.Viper
	Logger     *logrus.Logger
	Key        []byte
	Report     ReportInfo
	Status     AppStatus
	UCSPM      UCSPMInfo
	UCS        UCSInfo
	Results    []CombinedResults
	Action     string
	Version    string
	Commands   []CommandInfo
}

type UCSPMDeviceInfo struct {
	uid               string
	uuid              string
	ignore            bool
	name              string
	model             string
	ishypervisor      bool
	hypervisorName    string
	hypervisorVersion string
	ucspmName         string
	hasHypervisor     bool
}

type ReportData struct {
	timestamp string
	value     float64
}

type dataSlice []ReportData

// Len is part of sort.Interface.
func (d dataSlice) Len() int {
	return len(d)
}

// Swap is part of sort.Interface.
func (d dataSlice) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}

// Less is part of sort.Interface. We use count as the value to sort by
func (d dataSlice) Less(i, j int) bool {
	return d[i].timestamp < d[j].timestamp
}
