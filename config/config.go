package config

import (
	"log"

	"github.com/BurntSushi/toml"
)

//AppConfig App configuration
type AppConfig struct {
	Environment string
	Log         map[string]Log
}

//Log configuration
type Log struct {
	Tracelevel  string
	Stacktrace  bool
	Erroroutput bool
	Caller      bool
	CallerSkip  int `toml:"caller_skip"`
}

var apps map[string]AppConfig

//Load toml file
func Load(filepath string) {
	if _, err := toml.DecodeFile(filepath, &apps); err != nil {
		log.Fatal("ERR: ", err)
	}
}

//GetApp returns AppConfig for particular app
func GetApp(a string) *AppConfig {
	app, ok := apps[a]

	if !ok {
		return &AppConfig{}
	}
	return &app
}

//GetLog Log config from certain app environment
func GetLog(env string) *map[string]Log {
	app, ok := apps[env]

	if !ok {
		return nil
	}
	return &app.Log
}
