package resource

import (
	"AynaLivePlayer/config"
	"io/ioutil"
)

var ProgramIcon = []byte{}
var EmptyImage = []byte{}

func init() {
	loadResource(config.GetAssetPath("icon.jpg"), &ProgramIcon)
	loadResource(config.GetAssetPath("empty.png"), &EmptyImage)
}

func loadResource(path string, res *[]byte) {
	if file, err := ioutil.ReadFile(path); err == nil {
		*res = file
	}
}
