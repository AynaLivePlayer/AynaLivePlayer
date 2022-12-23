package model

import "fmt"

type LiveRoom struct {
	ClientName  string
	ID          string
	AutoConnect bool
}

func (r *LiveRoom) String() string {
	return fmt.Sprintf("<LiveRooms %s:%s>", r.ClientName, r.ID)
}

func (r *LiveRoom) Title() string {
	return fmt.Sprintf("%s-%s", r.ClientName, r.ID)
}
