package adapter

import "AynaLivePlayer/core/model"

type IApplication interface {
	Version() model.VersionInfo
	LatestVersion() model.VersionInfo
	CheckUpdate() error
}
