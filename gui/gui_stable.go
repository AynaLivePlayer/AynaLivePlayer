//go:build nosource

package gui

import (
	"AynaLivePlayer/core/model"
	"AynaLivePlayer/pkg/config"
	"fmt"
)

func getAppTitle() string {
	return fmt.Sprintf("%s Ver %s - 测试版 仅供开发人员测试使用 请勿用于其他用途", config.ProgramName, model.Version(config.Version))
}
