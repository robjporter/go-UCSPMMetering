package flags

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"

	functions "github.com/robjporter/go-functions"
	"github.com/robjporter/go-functions/as"
	"github.com/robjporter/go-functions/kingpin"
	"github.com/robjporter/go-functions/viper"
	"github.com/robjporter/go-functions/yaml"
)

var (
	add    = kingpin.Command("add", "Register a new UCS domain.")
	update = kingpin.Command("update", "Update a UCS domain.")
	delete = kingpin.Command("delete", "Remove a UCS domain.")
	show   = kingpin.Command("show", "Show a UCS domain.")
	run    = kingpin.Command("run", "Run the main application.")
	output = kingpin.Command("output", "Configure the output file.")
	input  = kingpin.Command("input", "Configure the input file.")

	debug = kingpin.Flag("debug", "Enable debug mode.").Bool()

	addUCS    = add.Command("ucs", "Add a UCS Domain")
	updateUCS = update.Command("ucs", "Update a UCS Domain")
	deleteUCS = delete.Command("ucs", "Delete a UCS Domain")
	showUCS   = show.Command("ucs", "Show a UCS Domain")

	addUCSPM    = add.Command("ucspm", "Add a UCSPM Domain")
	updateUCSPM = update.Command("ucspm", "Update a UCSPM Domain")
	deleteUCSPM = delete.Command("ucspm", "Delete a UCSPM Domain")
	showUCSPM   = show.Command("ucspm", "Show a UCS Performance Manager")

	addUCSIP       = addUCS.Flag("ip", "IP Address or DNS name for UCS Manager, without http(s).").Required().IP()
	addUCSUsername = addUCS.Flag("username", "Name of user.").Required().String()
	addUCSPassword = addUCS.Flag("password", "Password for user in plain text.").Required().String()

	updateUCSIP       = updateUCS.Flag("ip", "IP Address or DNS name for UCS Manager, without http(s).").Required().IP()
	updateUCSUsername = updateUCS.Flag("username", "Name of user.").Required().String()
	updateUCSPassword = updateUCS.Flag("password", "Password for user in plain text.").Required().String()

	deleteUCSIP = deleteUCS.Flag("ip", "IP Address or DNS name for UCS Manager, without http(s).").Required().IP()

	showUCSIP = showUCS.Flag("ip", "IP Address or DNS name for UCS Manager, without http(s).").Required().IP()

	addUCSPMIP       = addUCSPM.Flag("ip", "IP Address or DNS name for UCS Performance Manager, without http(s).").Required().IP()
	addUCSPMUsername = addUCSPM.Flag("username", "Name of user.").Required().String()
	addUCSPMPassword = addUCSPM.Flag("password", "Password for user in plain text.").Required().String()

	updateUCSPMIP       = updateUCSPM.Flag("ip", "IP Address or DNS name for UCS Performance Manager, without http(s).").Required().IP()
	updateUCSPMUsername = updateUCSPM.Flag("username", "Name of user.").Required().String()
	updateUCSPMPassword = updateUCSPM.Flag("password", "Password for user in plain text.").Required().String()

	outputFile = output.Flag("set", "Configure the output filename, where the UUID and serial numbers will be saved.").Required().String()
	inputFile  = input.Flag("set", "Configure the input filename, where the UUID will be read from.").Required().String()

	ConfigFile = ""
	systems    []System
	key        = []byte("CiscoFinanceOpenPay12345")
)

func LoadConfig(filename string) {
	configName := ""
	configExtension := ""
	configPath := ""

	splits := strings.Split(filepath.Base(filename), ".")
	if len(splits) == 2 {
		configName = splits[0]
		configExtension = splits[1]
	}
	configPath = filepath.Dir(filename)

	viper.SetConfigName(configName)
	viper.SetConfigType(configExtension)
	viper.AddConfigPath(configPath)

	createBlankConfig(filename)
	ConfigFile = filename
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}

func createBlankConfig(filename string) {
	if !functions.Exists(filename) {
		viper.Set("eula.agreed", false)
		viper.Set("input.file", "uuid.json")
		viper.Set("output.file", "output.csv")
		saveConfig()
	}
}

func saveConfig() {
	if len(systems) > 0 {
		items := processSystems()
		viper.Set("ucs.systems", items)
	}
	out, err := yaml.Marshal(viper.AllSettings())
	if err == nil {
		fp, err := os.Create(ConfigFile)
		if err == nil {
			defer fp.Close()
			_, err = fp.Write(out)
		}
	}
}

func processSystems() []interface{} {
	var items []interface{}
	var item map[string]interface{}
	for i := 0; i < len(systems); i++ {

		item = make(map[string]interface{})
		item["url"] = systems[i].ip
		item["username"] = systems[i].username
		item["password"] = systems[i].password
		items = append(items, item)
	}
	return items
}

func setInputFilename(filename string) {
	viper.Set("input.file", filename)
}

func setOutputFilename(filename string) {
	viper.Set("output.file", filename)
}

func deleteUCSPMSystem() bool {
	getAllSystems()
	if viper.IsSet("ucspm.url") {
		viper.Set("ucspm", false)
		return true
	}
	return false
}

func showUCSPMSystem() {
	if viper.IsSet("ucspm.url") {
		fmt.Println("UCS Performance Manager URL:     > ", viper.GetString("ucspm.url"))
		fmt.Println("UCS Performance Manager Username:> ", viper.GetString("ucspm.username"))
		fmt.Println("UCS Performance Manager Password:> ", viper.GetString("ucspm.password"))
	}
}

func checkUCSPMExists(UCSIP net.IP) bool {
	getAllSystems()
	if viper.IsSet("ucspm.url") {
		return true
	}
	return false
}

func addUCSPMSystem(UCSIP net.IP, username, password string) bool {
	if as.ToString(UCSIP) != "" {
		if as.ToString(username) != "" {
			if as.ToString(password) != "" {
				viper.Set("ucspm.url", as.ToString(UCSIP))
				viper.Set("ucspm.username", as.ToString(username))
				viper.Set("ucspm.password", encryptPassword(as.ToString(password)))
				return true
			} else {
				fmt.Println("UCS Performance Manager password cannot be blank.")
			}
		} else {
			fmt.Println("UCS Performance Manager username cannot be blank.")
		}
	} else {
		fmt.Println("UCS Performance Manager URL cannot be blank.")
	}
	return false
}

func updateUCSPMSystem(UCSIP net.IP, addUCSUsername, addUCSPassword string) bool {
	return addUCSPMSystem(UCSIP, addUCSUsername, addUCSPassword)
}

func displayUCSSystem(UCSIP net.IP) {
	for i := 0; i < len(systems); i++ {
		if systems[i].ip == as.ToString(UCSIP) {
			fmt.Println("UCS Domain Url:     > ", systems[i].ip)
			fmt.Println("UCS Domain Username:> ", systems[i].username)
			fmt.Println("UCS Domain Password:> ", systems[i].password)
		}
	}
}

func deleteUCSSystem(UCSIP net.IP) bool {
	for i := 0; i < len(systems); i++ {
		if systems[i].ip == as.ToString(UCSIP) {
			systems = append(systems[:i], systems[i+1:]...)
		}
	}
	return true
}

func updateUCSSystem(UCSIP net.IP, addUCSUsername, addUCSPassword string) bool {
	for i := 0; i < len(systems); i++ {
		if systems[i].ip == as.ToString(UCSIP) {
			systems[i].username = addUCSUsername
			systems[i].password = encryptPassword(addUCSPassword)
		}
	}
	return true
}

func checkUCSExists(ip net.IP) bool {
	if viper.IsSet("ucs.systems") {
		getAllSystems()
		for i := 0; i < len(systems); i++ {
			sysip := ip.String()
			if strings.TrimSpace(systems[i].ip) == strings.TrimSpace(sysip) {
				return true
			}
		}
		return false
	}
	return false
}

func getAllSystems() {
	tmp := as.ToSlice(viper.Get("ucs.systems"))
	readSystems(tmp)
}

func encryptPassword(password string) string {
	return functions.Encrypt(key, []byte(password))
}

func decryptPassword(password string) string {
	return functions.Decrypt(key, password)
}

func addUCSSystem(addUCSIP net.IP, addUCSUsername, addUCSPassword string) bool {
	tmp := System{}
	tmp.ip = addUCSIP.String()
	tmp.username = addUCSUsername
	tmp.password = encryptPassword(addUCSPassword)
	systems = append(systems, tmp)
	return true
}

func readSystems(ucss []interface{}) bool {
	for i := 0; i < len(ucss); i++ {
		var newlist map[string]string
		newlist = as.ToStringMapString(ucss[i])
		tmp := System{}
		tmp.ip = newlist["url"]
		tmp.username = newlist["username"]
		tmp.password = newlist["password"]
		systems = append(systems, tmp)
	}

	return true
}

func ProcessCommandLineArguments() string {
	switch kingpin.Parse() {
	case "run":
		return "RUN"
	case "add ucs":
		if !checkUCSExists(*addUCSIP) {
			if addUCSSystem(*addUCSIP, *addUCSUsername, *addUCSPassword) {
				saveConfig()
				fmt.Println("New system has been added successfully")
			} else {
				fmt.Println("System could not be added.")
			}
		} else {
			fmt.Println("UCS System already exsists.")
		}
	case "update ucs":
		if checkUCSExists(*updateUCSIP) {
			if updateUCSSystem(*updateUCSIP, *updateUCSUsername, *updateUCSPassword) {
				saveConfig()
				fmt.Println("UCS System ", as.ToString(*updateUCSIP), " has been updated successfully")
			}
		} else {
			println("The UCS Domain: ", as.ToString(*updateUCSIP), " is not in the config file, please try running add.")
		}
	case "delete ucs":
		if checkUCSExists(*deleteUCSIP) {
			if deleteUCSSystem(*deleteUCSIP) {
				saveConfig()
				fmt.Println("UCS System ", as.ToString(*deleteUCSIP), " has been deleted successfully")
			}
		} else {
			println("The UCS Domain: ", as.ToString(*deleteUCSIP), " is not in the config file, please try running add.")
		}
	case "show ucs":
		if checkUCSExists(*showUCSIP) {
			displayUCSSystem(*showUCSIP)
		} else {
			println("The UCS Domain: ", as.ToString(*showUCSIP), " is not in the config file, please try running add.")
		}
	case "add ucspm":
		if !checkUCSPMExists(*addUCSPMIP) {
			if addUCSPMSystem(*addUCSPMIP, *addUCSPMUsername, *addUCSPMPassword) {
				saveConfig()
				fmt.Println("UCS Performance Manager system ", as.ToString(*addUCSPMIP), " has been added successfully")
			}
		} else {
			fmt.Println("A UCS Performance Manager instance already exists in the config file, please view, remove or update.")
		}
	case "update ucspm":
		if checkUCSPMExists(*updateUCSPMIP) {
			if updateUCSPMSystem(*updateUCSPMIP, *updateUCSPMUsername, *updateUCSPMPassword) {
				saveConfig()
				fmt.Println("UCS Performance Manager system ", as.ToString(*updateUCSPMIP), " has been updated successfully")
			}
		}
	case "delete ucspm":
		if deleteUCSPMSystem() {
			saveConfig()
			fmt.Println("UCS Performance Manager system has been deleted successfully")
		} else {
			fmt.Println("A UCS Performance Manager instance does not exist in the config file, please add one.")
		}
	case "show ucspm":
		showUCSPMSystem()
	case "input":
		setInputFilename(*inputFile)
		saveConfig()
		fmt.Println("Application input file has been updated successfully")
	case "output":
		setOutputFilename(*outputFile)
		saveConfig()
		fmt.Println("Application output file has been updated  successfully")
	}
	return ""
}
