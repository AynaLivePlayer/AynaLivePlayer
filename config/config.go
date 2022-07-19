package config

import (
	"fmt"
	"gopkg.in/ini.v1"
	"path"
)

const (
	ProgramName = "卡西米尔唱片机"
	Version     = "alpha 0.8.6"
)

const (
	ConfigPath = "./config.ini"
	AssetsPath = "./assets"
)

func GetAssetPath(name string) string {
	return path.Join(AssetsPath, name)
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
	ConfigFile, err = ini.Load(ConfigPath)
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
