package model

type User struct {
	Name string
}

func ApplyUser(medias []*Media, user interface{}) {
	for _, m := range medias {
		m.User = user
	}
}

func ToSpMedia(media *Media, user *User) *Media {
	media = media.Copy()
	media.User = user
	return media
}
