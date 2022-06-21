package config

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
	"path"
)

const VERSION = "alpha 0.4"

const CONFIG_PATH = "./config.ini"
const Assests_PATH = "./assets"

func GetAssetPath(name string) string {
	return path.Join(Assests_PATH, name)
}

type LogConfig struct {
	Path  string
	Level logrus.Level
}

var Log = &LogConfig{
	Path:  "./log.txt",
	Level: logrus.InfoLevel,
}

type LiveRoomConfig struct {
	History []string
}

var LiveRoom = &LiveRoomConfig{History: []string{"9076804", "3819533"}}

type PlayerConfig struct {
	Playlists         []string
	PlaylistsProvider []string
	PlaylistIndex     int
	PlaylistRandom    bool
}

var Player = &PlayerConfig{
	Playlists:         []string{"116746576", "646548465"},
	PlaylistsProvider: []string{"netease", "netease"},
	PlaylistIndex:     0,
	PlaylistRandom:    true,
}

type ProviderConfig struct {
	Priority []string
	LocalDir string
}

var Provider = &ProviderConfig{
	Priority: []string{"local", "netease", "kuwo", "bilibili"},
	LocalDir: "./music",
}

func init() {
	cfg, err := ini.Load(CONFIG_PATH)
	if err != nil {
		fmt.Println("config not found")
		SaveToConfigFile(CONFIG_PATH)
		return
	}
	err = cfg.Section("Log").MapTo(Log)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = cfg.Section("LiveRoom").MapTo(LiveRoom)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(cfg.Section("Player").GetKey("Playlists"))
	err = cfg.Section("Player").MapTo(Player)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = cfg.Section("Provider").MapTo(Provider)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func SaveToConfigFile(filename string) error {
	cfg := ini.Empty()
	err := ini.ReflectFrom(cfg, &struct {
		Log      *LogConfig
		LiveRoom *LiveRoomConfig
		Player   *PlayerConfig
		Provider *ProviderConfig
	}{
		Log:      Log,
		LiveRoom: LiveRoom,
		Player:   Player,
		Provider: Provider,
	})
	if err != nil {
		return err
	}
	return cfg.SaveTo(filename)
}
