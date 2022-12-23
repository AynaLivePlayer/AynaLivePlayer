package config

import (
	"fmt"
	"gopkg.in/ini.v1"
	"path"
)

const (
	ProgramName = "卡西米尔唱片机"
	Version     = "beta 0.9.5"
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
	OnLoad()
	OnSave()
}

type BaseConfig struct {
}

func (c *BaseConfig) OnLoad() {
}

func (c *BaseConfig) OnSave() {
}

var ConfigFile *ini.File
var Configs = make([]Config, 0)

func LoadConfig(cfg Config) {
	sec, err := ConfigFile.GetSection(cfg.Name())
	if err == nil {
		_ = sec.MapTo(cfg)
	}
	cfg.OnLoad()
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
	for _, cfg := range []Config{Log, General} {
		LoadConfig(cfg)
	}
}

func SaveToConfigFile(filename string) error {
	cfgFile := ini.Empty()
	for _, cfg := range Configs {
		cfg.OnSave()
		if err := cfgFile.Section(cfg.Name()).ReflectFrom(cfg); err != nil {
			fmt.Println(err)
		}
	}
	return cfgFile.SaveTo(filename)
}
