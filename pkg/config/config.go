package config

import (
	"encoding/json"
	"fmt"
	"gopkg.in/ini.v1"
	"os"
	"path"
)

const (
	ProgramName        = "卡西米尔唱片机"
	Version     uint32 = 0x010005
)

const (
	ConfigPath = "./config.ini"
	AssetsPath = "./assets"
)

func GetAssetPath(name string) string {
	return path.Join(AssetsPath, name)
}

var DEFAULT_CONFIGS = []Config{General}

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
var Environment = CurrentEnvironment()

func LoadFromFile(filename string) {
	var err error
	ConfigFile, err = ini.Load(filename)
	if err != nil {
		fmt.Println("config not found, using default config")
		ConfigFile = ini.Empty()
	}
	for _, cfg := range DEFAULT_CONFIGS {
		LoadConfig(cfg)
	}
}

func init() {
	ConfigFile = ini.Empty()
	for _, cfg := range DEFAULT_CONFIGS {
		LoadConfig(cfg)
	}
}

func LoadConfig(cfg Config) {
	sec, err := ConfigFile.GetSection(cfg.Name())
	if err == nil {
		_ = sec.MapTo(cfg)
	}
	cfg.OnLoad()
	Configs = append(Configs, cfg)
	return
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

func LoadJson(path string, dst any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dst)
}

func SaveJson(path string, dst any) error {
	data, err := json.MarshalIndent(dst, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0666)
}
