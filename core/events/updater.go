package events

import "AynaLivePlayer/core/model"

const CheckUpdateCmd = "cmd.update.check"

type CheckUpdateCmdEvent struct {
}

const CheckUpdateResultUpdate = "update.update.check"

type CheckUpdateResultUpdateEvent struct {
	HasUpdate bool
	Info      model.VersionInfo
}
