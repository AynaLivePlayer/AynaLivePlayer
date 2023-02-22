package internal

import (
	"AynaLivePlayer/common/config"
	"AynaLivePlayer/core/model"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/tidwall/gjson"
)

type AppBilibiliChannel struct {
	latestVersion model.Version
}

func (app *AppBilibiliChannel) Version() model.VersionInfo {
	return model.VersionInfo{
		model.Version(config.Version), "",
	}
}

func (app *AppBilibiliChannel) LatestVersion() model.VersionInfo {
	return model.VersionInfo{
		app.latestVersion,
		fmt.Sprintf("v%s\n\n[https://play-live.bilibili.com/details/1661006726438](https://play-live.bilibili.com/details/1661006726438)", app.latestVersion),
	}
}

func (app *AppBilibiliChannel) CheckUpdate() error {
	uri := "https://api.live.bilibili.com/xlive/virtual-interface/v1/app/detail?app_id=1661006726438"
	resp, err := resty.New().R().Get(uri)
	if err != nil {
		return err
	}
	lv := model.VersionFromString(gjson.ParseBytes(resp.Body()).Get("data.version").String())
	if lv == 0 {
		return errors.New("failed to get latest version")
	}
	app.latestVersion = lv
	return nil
}
