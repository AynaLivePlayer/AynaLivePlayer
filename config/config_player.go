package config

type _PlayerConfig struct {
	PlaylistData string
	Playlists    []*PlayerPlaylist `ini:"-"`
	//PlaylistsProvider  []string
	PlaylistIndex      int
	PlaylistRandom     bool
	UserPlaylistRandom bool
	AudioDevice        string
	Volume             float64
	SkipPlaylist       bool
}

type PlayerPlaylist struct {
	ID       string
	Provider string
}

func (c *_PlayerConfig) Name() string {
	return "Player"
}

func (c *_PlayerConfig) OnLoad() {
	//c.Playlists = make([]*PlayerPlaylist, 0)
	_ = LoadJson(c.PlaylistData, &c.Playlists)
}

func (c *_PlayerConfig) OnSave() {
	_ = SaveJson(c.PlaylistData, &c.Playlists)
}

var Player = &_PlayerConfig{
	PlaylistData: "playlists.json",
	Playlists: []*PlayerPlaylist{
		{
			"2382819181",
			"netease",
		},
		{
			"4987059624",
			"netease",
		},
		{
			"list1",
			"local",
		},
	},
	PlaylistIndex:      0,
	PlaylistRandom:     true,
	UserPlaylistRandom: false,
	AudioDevice:        "auto",
	Volume:             100,
	SkipPlaylist:       false,
}
