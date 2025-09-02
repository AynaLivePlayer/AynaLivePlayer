package model

import "github.com/AynaLivePlayer/liveroom-sdk"

type LiveRoomConfig struct {
	AutoConnect bool `json:"auto_connect"`
}

type LiveRoom struct {
	LiveRoom liveroom.LiveRoom `json:"live_room"`
	Config   LiveRoomConfig    `json:"config"`
	Title    string            `json:"title"`
	Status   bool              `json:"status"`
}

func (r *LiveRoom) DisplayName() string {
	if r.Title != "" {
		return r.Title
	}
	return r.LiveRoom.Identifier()
}

type LiveRoomProviderInfo struct {
	Name        string
	Description string
}
