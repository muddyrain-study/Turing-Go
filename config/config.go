package config

import (
	"github.com/Unknwon/goconfig"
	"log"
	"os"
)

var File *goconfig.ConfigFile

const confFile = "/conf/conf.ini"

func init() {
	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	configPath := currentDir + confFile
	if !fileExist(configPath) {
		panic("Config file not found: " + configPath)
	}
	len := len(os.Args)
	if len > 1 {
		dir := os.Args[1]
		if dir != "" {
			configPath = dir + confFile
		}
	}
	File, err = goconfig.LoadConfigFile(configPath)
	if err != nil {
		log.Fatal("Error loading configuration file: " + err.Error())
	}
}
func fileExist(fileName string) bool {
	_, err := os.Stat(fileName)
	return err == nil || os.IsExist(err)
}
