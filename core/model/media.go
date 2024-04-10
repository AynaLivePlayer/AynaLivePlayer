package model

import (
	"github.com/AynaLivePlayer/liveroom-sdk"
	"github.com/AynaLivePlayer/miaosic"
)

type User struct {
	Name string
}

var PlaylistUser = User{Name: "Playlists"}
var SystemUser = User{Name: "System"}
var HistoryUser = User{Name: "History"}

type Media struct {
	Info miaosic.MediaInfo
	User interface{}
}

func (m *Media) IsLiveRoomUser() bool {
	_, ok := m.User.(liveroom.User)
	return ok
}

func (m *Media) ToUser() User {
	if u, ok := m.User.(User); ok {
		return u
	}
	return User{Name: m.DanmuUser().Username}
}

func (m *Media) DanmuUser() liveroom.User {
	if u, ok := m.User.(liveroom.User); ok {
		return u
	}
	return liveroom.User{}
}
