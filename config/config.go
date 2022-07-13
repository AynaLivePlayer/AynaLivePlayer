package config

import (
	"fmt"
	"gopkg.in/ini.v1"
	"path"
)

const VERSION = "alpha 0.8.3"

const CONFIG_PATH = "./config.ini"
const Assests_PATH = "./assets"

func GetAssetPath(name string) string {
	return path.Join(Assests_PATH, name)
}

type Config interface {
	Name() string
}

var ConfigFile *ini.File
var Configs = make([]Config, 0)

func LoadConfig(cfg Config) {
	sec, err := ConfigFile.GetSection(cfg.Name())
	if err == nil {
		_ = sec.MapTo(cfg)
	}
	Configs = append(Configs, cfg)
	return
}

func init() {
	var err error
	ConfigFile, err = ini.Load(CONFIG_PATH)
	if err != nil {
		fmt.Println("config not found, using default config")
		ConfigFile = ini.Empty()
	}
	for _, cfg := range []Config{Log, LiveRoom, Player, Provider, General} {
		LoadConfig(cfg)
	}
}

func SaveToConfigFile(filename string) error {
	cfgFile := ini.Empty()
	for _, cfg := range Configs {
		if err := cfgFile.Section(cfg.Name()).ReflectFrom(cfg); err != nil {
			fmt.Println(err)
		}
	}
	return cfgFile.SaveTo(filename)
}
