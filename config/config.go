package config

import (
	"log"

	"github.com/BurntSushi/toml"
)

//AppConfig ...
type AppConfig struct {
	Environment string
	Log         map[string]Log
}

//Log ...
type Log struct {
	Tracelevel  string //= ""
	Stacktrace  bool   //=  false
	Erroroutput bool   //= false
	Caller      bool   //= false
	CallerSkip  int    `toml:"caller_skip"`
}

var apps map[string]AppConfig

//Load toml file
func Load(filepath string) {
	if _, err := toml.DecodeFile(filepath, &apps); err != nil { //.Decode(tomlData, &conf); err != nil {
		log.Fatal("ERR: ", err)
	}
	log.Println(apps)
}

//GetApp returns AppConfig for particular app
func GetApp(a string) *AppConfig {
	app, ok := apps[a]

	if !ok {
		return &AppConfig{}
	}
	return &app
}

//GetLog ...
func GetLog(env string) *map[string]Log {
	app, ok := apps[env]

	if !ok {
		return nil
	}
	return &app.Log
	// appLog, found := app.Log[s]
	//
	// if !found {
	// 	return nil
	// }
	// return &appLog
}
