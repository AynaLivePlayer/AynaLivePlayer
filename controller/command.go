package controller

import (
	"AynaLivePlayer/event"
	"AynaLivePlayer/liveclient"
	"strings"
)

var Commands []DanmuCommandExecutor

type DanmuCommandExecutor interface {
	Match(command string) bool
	Execute(command string, args []string, danmu *liveclient.DanmuMessage)
}

func AddCommand(executors ...DanmuCommandExecutor) {
	Commands = append(Commands, executors...)
}

func danmuCommandHandler(event *event.Event) {
	danmu := event.Data.(*liveclient.DanmuMessage)
	args := strings.Split(danmu.Message, " ")
	if len(args[0]) == 0 {
		return
	}
	for _, cmd := range Commands {
		if cmd.Match(args[0]) {
			cmd.Execute(args[0], args[1:], danmu)
		}
	}
}
