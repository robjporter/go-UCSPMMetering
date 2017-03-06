package ucspm

type device struct {
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

type Application struct {
	configFile string
	routers    map[string]string
	devices    []device
	DEBUG      bool
	tidCount   int
	host       string
	username   string
	password   string
	outputFile string
}
