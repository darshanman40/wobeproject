package config

import (
	"log"

	"github.com/BurntSushi/toml"
)

//AppConfig ...
type AppConfig struct {
	Enabled  bool `toml:"Enabled"`
	Port     string
	URLPath  string
	Hosts    []string
	Encoding string
	LoadGen  LoadGen `toml:"loadgen"`
}

//LoadGen ...
type LoadGen struct {
	TPS      int64  `toml:"tps"`
	Duration string `toml:"duration"`
}

var apps map[string]AppConfig

//Load toml file
func Load(env string) {
	if _, err := toml.DecodeFile("config/"+env+".toml", &apps); err != nil { //.Decode(tomlData, &conf); err != nil {
		// handle error
		log.Fatal("ERR: ", err)
	}
	log.Println(apps)
}

//GetApp returns AppConfig for particular app
func GetApp(a string) *AppConfig {
	app, ok := apps[a]

	if !ok {
		return &AppConfig{} //, errors.New("Config not found")
	}
	return &app
}
