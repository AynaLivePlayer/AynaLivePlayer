package controller

import (
	"AynaLivePlayer/liveclient"
	"strings"
)

var CommandDiange = &Diange{}

type Diange struct {
}

func (d Diange) Match(command string) bool {
	for _, c := range []string{"点歌"} {
		if command == c {
			return true
		}
	}
	return false
}

func (d Diange) Execute(command string, args []string, danmu *liveclient.DanmuMessage) {
	keyword := strings.Join(args, " ")
	Add(keyword, &danmu.User)
}
