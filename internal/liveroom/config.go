package liveroom

import (
	"AynaLivePlayer/core/model"
	"AynaLivePlayer/pkg/config"
)

type _cfg struct {
	ApiServer    string
	LiveRoomPath string
	liveRooms    []model.LiveRoom
}

func (c *_cfg) Name() string {
	return "LiveRoom"
}

func (c *_cfg) OnLoad() {
	_ = config.LoadJson(c.LiveRoomPath, &c.liveRooms)
}

func (c *_cfg) OnSave() {
	_ = config.SaveJson(c.LiveRoomPath, &c.liveRooms)
}

var cfg = &_cfg{
	ApiServer:    "http://localhost:9090",
	LiveRoomPath: "liveroom.json",
	liveRooms:    make([]model.LiveRoom, 0),
}
