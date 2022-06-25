package config

type _PlayerConfig struct {
	Playlists         []string
	PlaylistsProvider []string
	PlaylistIndex     int
	PlaylistRandom    bool
}

func (c *_PlayerConfig) Name() string {
	return "Player"
}

var Player = &_PlayerConfig{
	Playlists:         []string{"2382819181", "116746576", "646548465"},
	PlaylistsProvider: []string{"netease", "netease", "netease"},
	PlaylistIndex:     0,
	PlaylistRandom:    true,
}
