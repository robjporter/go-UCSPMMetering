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

type UCSPMInfo struct {
	Routers  map[string]string
	TidCount int
	Devices  []UCSPMDeviceInfo
	host     string
	username string
	password string
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

type Application struct {
	ConfigFile string
	Debug      bool
	Config     *viper.Viper
	UCSSystems []UCSSystemInfo
	Logger     *logrus.Logger
	Key        []byte
	Report     ReportInfo
	Status     AppStatus
	UCSPM      UCSPMInfo
}

type UCSPMDeviceInfo struct {
	uid               string
	uuid              string
	ignore            bool
	name              string
	model             string
	hypervisor        bool
	hypervisorName    string
	hypervisorVersion string
	ucspmName         string
	hasHypervisor     bool
}
