package player

import (
	"fmt"
	"strconv"
	"testing"
)

func TestPlaylist_Insert(t *testing.T) {
	pl := NewPlaylist("asdf", PlaylistConfig{RandomNext: false})
	for i := 0; i < 10; i++ {
		pl.Insert(-1, &Media{Url: strconv.Itoa(i)})
	}
	pl.Insert(3, &Media{Url: "a"})
	pl.Insert(0, &Media{Url: "b"})
	pl.Insert(-2, &Media{Url: "x"})
	pl.Insert(-1, &Media{Url: "h"})
	for i := 0; i < pl.Size(); i++ {
		fmt.Print(pl.Playlist[i].Url, " ")
	}

}
