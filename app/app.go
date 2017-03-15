package app

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"../flags"
	"../functions"

	functions2 "github.com/robjporter/go-functions"
	"github.com/robjporter/go-functions/as"
	"github.com/robjporter/go-functions/banner"
	"github.com/robjporter/go-functions/colors"
	"github.com/robjporter/go-functions/lfshook"
	"github.com/robjporter/go-functions/logrus"
	"github.com/robjporter/go-functions/terminal"
	"github.com/robjporter/go-functions/timing"
	"github.com/robjporter/go-functions/viper"
	"github.com/robjporter/go-functions/yaml"
)

//TODO: UNMATCHED.json missing!
//TODO: File has been saved successfully.             Filename=unmatched.json happens as 5th line when starting app!

var (
	Core Application
)

func (a *Application) createBlankConfig(filename string) {
	if !functions2.Exists(filename) {
		a.LogInfo("Creating a new default configuration file.", nil, true)
		a.Config.Set("eula.agreed", false)
		a.Config.Set("output.file", "output.csv")
		a.Config.Set("debug", false)
		a.Config.Set("metrics.run", 0)
		a.Config.Set("metrics.clean", 0)
		a.Config.Set("metrics.adducs", 0)
		a.Config.Set("metrics.updateucs", 0)
		a.Config.Set("metrics.deleteucs", 0)
		a.Config.Set("metrics.showucs", 0)
		a.Config.Set("metrics.showall", 0)
		a.Config.Set("metrics.adducspm", 0)
		a.Config.Set("metrics.updateucspm", 0)
		a.Config.Set("metrics.deleteucspm", 0)
		a.Config.Set("metrics.showucspm", 0)
		a.Config.Set("metrics.setinput", 0)
		a.Config.Set("metrics.setoutput", 0)
		a.saveConfig()
	}
}

func (a *Application) Start() {
	timing.Timer("CORE")
}

func (a *Application) Finish() {
	if a.Action == "CLEAN" {
		fmt.Println("Application run finished. Timer =", timing.Timer("CORE"))
	} else {
		a.LogInfo("Application run finished.", map[string]interface{}{"Timer": timing.Timer("CORE")}, false)
	}
}

func (a *Application) EncryptPassword(password string) string {
	return functions2.Encrypt(a.Key, []byte(password))
}

func (a *Application) DecryptPassword(password string) string {
	return functions2.Decrypt(a.Key, password)
}

func (a *Application) getReportDates(month, year string) (string, string) {
	if month == "" {
		month = functions.CurrentMonthName()
	} else {
		tmp := functions.IsMonth(month)
		if tmp != "" {
			month = tmp
		} else {
			month = functions.CurrentMonthName()
		}
	}
	if year == "" {
		year = functions.CurrentYear()
	} else {
		tmp := functions.IsYear(year)
		if tmp != "" {
			year = tmp
		} else {
			year = functions.CurrentYear()
		}
	}
	return month, year
}
func (a *Application) init() {
	a.Config = viper.New()
	a.Logger = logrus.New()
	a.Logger.Level = logrus.DebugLevel
	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "02-01-2006 15:04:05.000"
	customFormatter.FullTimestamp = true
	a.Logger.Formatter = customFormatter
	a.Logger.Out = os.Stdout
	ts := as.ToString(time.Now().Unix())
	a.Logger.Hooks.Add(lfshook.NewHook(lfshook.PathMap{
		logrus.InfoLevel:  "data/info-" + ts + ".log",
		logrus.ErrorLevel: "data/error-" + ts + ".log",
		logrus.WarnLevel:  "data/warn-" + ts + ".log",
		logrus.DebugLevel: "data/debug-" + ts + ".log",
		logrus.FatalLevel: "data/fatal-" + ts + ".log",
	}))
	a.Key = []byte("CiscoFinanceOpenPay12345")
	a.displayBanner()
	os.Mkdir("data", 0700)
}

func (a *Application) displayBanner() {
	terminal.ClearScreen()
	banner.PrintNewFigure("UCS Metrics", "rounded", true)
	fmt.Println(colors.Color("Cisco Unified Computing System Metrics & Statistics Collection Tool v"+a.Version, colors.BRIGHTYELLOW))
	banner.BannerPrintLineS("=", 80)
}

func (a *Application) LoadConfig() {
	a.init()
	a.Log("Loading Configuration File.", nil, true)
	configName := ""
	configExtension := ""
	configPath := ""

	splits := strings.Split(filepath.Base(a.ConfigFile), ".")
	if len(splits) == 2 {
		configName = splits[0]
		configExtension = splits[1]
	}
	configPath = filepath.Dir(a.ConfigFile)

	a.Config.SetConfigName(configName)
	a.Config.SetConfigType(configExtension)
	a.Config.AddConfigPath(configPath)

	a.Log("Configuration File defined", map[string]interface{}{"Path": configPath, "Name": configName, "Extension": configExtension}, true)

	a.createBlankConfig(a.ConfigFile)

	err := a.Config.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
		os.Exit(0)
	}
	a.indexConfig()
	a.Log("Configuration File read successfully", nil, true)
}

func (a *Application) indexConfig() {
	a.Status.ucsCount = a.getAllUCSSystemsCount()
	a.Status.ucspmCount = a.getAllUCSPMSystemsCount()
	a.Status.eula = a.getEULAStatus()
	a.getAllSystems()
	a.Debug = a.Config.GetBool("debug")
}

func (a *Application) Log(message string, fields map[string]interface{}, debugMessage bool) {
	if debugMessage && a.Debug || !debugMessage {
		if fields != nil {
			a.Logger.WithFields(fields).Debug(message)
		} else {
			a.Logger.Debug(message)
		}
	}
}

func (a *Application) LogFatal(message string, fields map[string]interface{}) {
	if fields != nil {
		a.Logger.WithFields(fields).Fatal(message)
	} else {
		a.Logger.Fatal(message)
	}
}

func (a *Application) LogInfo(message string, fields map[string]interface{}, infoMessage bool) {
	if infoMessage && a.Debug || !infoMessage {
		if fields != nil {
			a.Logger.WithFields(fields).Info(message)
		} else {
			a.Logger.Info(message)
		}
	}
}

func (a *Application) LogWarn(message string, fields map[string]interface{}, warnMessage bool) {
	if warnMessage && a.Debug || !warnMessage {
		if fields != nil {
			a.Logger.WithFields(fields).Warn(message)
		} else {
			a.Logger.Warn(message)
		}
	}
}

func (a *Application) processSystems() []interface{} {
	var items []interface{}
	var item map[string]interface{}
	for i := 0; i < len(a.UCS.Systems); i++ {

		item = make(map[string]interface{})
		item["url"] = a.UCS.Systems[i].ip
		item["username"] = a.UCS.Systems[i].username
		item["password"] = a.UCS.Systems[i].password
		items = append(items, item)
	}
	return items
}

func (a *Application) Run() {
	a.LogInfo("Application", map[string]interface{}{"Version": a.Version}, false)
	a.LogInfo("Starting main application Run stage 1", nil, false)
	runtime.GOMAXPROCS(runtime.NumCPU())
	a.processResponse(flags.ProcessCommandLineArguments())
}

func (a *Application) saveConfig() {
	a.LogInfo("Saving configuration file.", nil, false)
	if len(a.UCS.Systems) > 0 {
		items := a.processSystems()
		a.Config.Set("ucs.systems", items)
	}
	out, err := yaml.Marshal(a.Config.AllSettings())
	if err == nil {
		fp, err := os.Create(a.ConfigFile)
		if err == nil {
			defer fp.Close()
			_, err = fp.Write(out)
		}
	}
	a.Log("Saving configuration file complete.", nil, true)
}

func (a *Application) saveFile(filename, data string) bool {
	filename = "data/" + filename
	ret := false
	f, err := os.Create(filename)
	if err == nil {
		_, err := f.Write([]byte(data))
		if err == nil {
			a.LogInfo("File has been saved successfully.", map[string]interface{}{"Filename": filename}, false)
			ret = true
		} else {
			a.LogInfo("There was a problem saving the file.", map[string]interface{}{"Error": err}, false)
		}
	}
	defer f.Close()
	return ret
}

func (a *Application) zipDataDir() {
	a.LogInfo("Preparing to archive output directory.", nil, false)
	functions2.Zipit("./data", "./Stage7-Complete-"+as.ToString(time.Now().Unix())+"-Data.zip")
	a.LogInfo("Archive created.", nil, false)
}

func (a *Application) addToCountMetrics(name string) {
	value := "metrics." + strings.ToLower(name)
	a.Config.Set(value, a.Config.GetInt(value)+1)
	a.saveConfig()
}
