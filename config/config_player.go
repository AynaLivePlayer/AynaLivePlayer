package config

type _PlayerConfig struct {
	Playlists         []string
	PlaylistsProvider []string
	PlaylistIndex     int
	PlaylistRandom    bool
	AudioDevice       string
	Volume            float64
	SkipPlaylist      bool
}

func (c *_PlayerConfig) Name() string {
	return "Player"
}

var Player = &_PlayerConfig{
	Playlists:         []string{"2382819181", "4987059624", "list1"},
	PlaylistsProvider: []string{"netease", "netease", "local"},
	PlaylistIndex:     0,
	PlaylistRandom:    true,
	AudioDevice:       "auto",
	Volume:            100,
	SkipPlaylist:      false,
}
