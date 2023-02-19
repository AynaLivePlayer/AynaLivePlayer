package model

import (
	"github.com/jinzhu/copier"
)

type Picture struct {
	Url  string
	Data []byte
}

func (p Picture) Exists() bool {
	return p.Url != "" || p.Data != nil
}

type Media struct {
	Title  string
	Artist string
	Cover  Picture
	Album  string
	Lyric  string
	Url    string
	Header map[string]string
	User   interface{}
	Meta   interface{}
}

func (m *Media) ToUser() *User {
	if u, ok := m.User.(*User); ok {
		return u
	}
	return &User{Name: m.DanmuUser().Username}
}

func (m *Media) DanmuUser() *DanmuUser {
	if u, ok := m.User.(*DanmuUser); ok {
		return u
	}
	return nil
}

func (m *Media) Copy() *Media {
	newMedia := &Media{}
	copier.Copy(newMedia, m)
	return newMedia
}
