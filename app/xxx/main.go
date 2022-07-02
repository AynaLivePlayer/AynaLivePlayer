package main

import (
	"AynaLivePlayer/plugin/textinfo"
)

func main() {
	x := &textinfo.TextInfo{}
	x.Enable()
	x.RenderTemplates()
}
