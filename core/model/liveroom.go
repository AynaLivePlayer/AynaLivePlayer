package model

import (
	"fmt"
)

type LiveRoom struct {
	ClientName    string
	ID            string
	Title         string
	AutoConnect   bool
	AutoReconnect bool
}

func (r *LiveRoom) String() string {
	return fmt.Sprintf("<LiveRooms %s:%s>", r.ClientName, r.ID)
}

func (r *LiveRoom) Identifier() string {
	return fmt.Sprintf("%s_%s", r.ClientName, r.ID)
}

type UserMedal struct {
	Name   string
	Level  int
	RoomID string
}

type DanmuUser struct {
	Uid       string
	Username  string
	Medal     UserMedal
	Admin     bool
	Privilege int
}

type DanmuMessage struct {
	User    DanmuUser
	Message string
}
