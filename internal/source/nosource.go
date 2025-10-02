//go:build nosource

package source

import (
	"github.com/AynaLivePlayer/miaosic"
)

func loadMediaProvider() {
	miaosic.RegisterProvider(&dummySource{})
}
