package updater

import (
	"AynaLivePlayer/core/events"
	"AynaLivePlayer/core/model"
	"AynaLivePlayer/global"
	"AynaLivePlayer/pkg/config"
	"AynaLivePlayer/pkg/event"
	"AynaLivePlayer/pkg/logger"
	"github.com/go-resty/resty/v2"
	"github.com/tidwall/gjson"
	"strconv"
)

var log logger.ILogger = nil

func Initialize() {
	log = global.Logger.WithPrefix("internal.updater")
	if config.General.AutoCheckUpdate {
		go func() {
			info, hasUpdate := CheckUpdate()
			if !hasUpdate {
				return
			}
			global.EventManager.CallA(
				events.CheckUpdateResultUpdate,
				events.CheckUpdateResultUpdateEvent{
					HasUpdate: hasUpdate,
					Info:      info,
				})
		}()
	}
	global.EventManager.RegisterA(
		events.CheckUpdateCmd, "internal.updater.handle",
		func(evt *event.Event) {
			info, hasUpdate := CheckUpdate()
			global.EventManager.CallA(
				events.CheckUpdateResultUpdate,
				events.CheckUpdateResultUpdateEvent{
					HasUpdate: hasUpdate,
					Info:      info,
				})
		})
}

func CheckUpdate() (model.VersionInfo, bool) {
	uri := config.General.InfoApiServer + "/api/version/check_update"
	resp, err := resty.New().R().SetQueryParam("client_version", strconv.Itoa(int(config.Version))).Get(uri)
	if err != nil {
		log.Errorf("failed to check update: %s", err.Error())
		return model.VersionInfo{}, false
	}
	result := gjson.ParseBytes(resp.Body())
	if !result.Get("data.has_update").Bool() {
		log.Infof("no update available")
		return model.VersionInfo{}, false
	}
	log.Infof("new version available: %s", model.Version(result.Get("data.latest.version").Uint()).String())
	return model.VersionInfo{
		Version: model.Version(result.Get("data.latest.version").Uint()),
		Info:    result.Get("data.latest.note").String(),
	}, true
}
