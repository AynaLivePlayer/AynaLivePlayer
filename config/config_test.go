package config

import (
	"fmt"
	"testing"
)

func TestCreate(t *testing.T) {
	fmt.Println(SaveToConfigFile(ConfigPath))
}

func TestLoad(t *testing.T) {
	fmt.Println(Log.Path)
	fmt.Println(Player.Playlists)
}
