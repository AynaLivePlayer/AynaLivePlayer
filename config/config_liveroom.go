package config

type _LiveRoomConfig struct {
	History []string
}

func (c *_LiveRoomConfig) Name() string {
	return "LiveRoom"
}

var LiveRoom = &_LiveRoomConfig{History: []string{"9076804", "3819533"}}
