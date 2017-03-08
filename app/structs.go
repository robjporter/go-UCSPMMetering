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

type ReportInfo struct {
	Month string
	Year  string
}

type Application struct {
	ConfigFile string
	Debug      bool
	Config     *viper.Viper
	UCSSystems []UCSSystemInfo
	Logger     *logrus.Logger
	Key        []byte
	Report     ReportInfo
}
