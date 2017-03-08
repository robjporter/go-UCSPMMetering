package ucs

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
