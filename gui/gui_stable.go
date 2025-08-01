//go:build nosource

package gui

import (
	"AynaLivePlayer/core/model"
	"AynaLivePlayer/pkg/config"
	"fmt"
)

func getAppTitle() string {
	return fmt.Sprintf("%s Ver %s - 正式版", config.ProgramName, model.Version(config.Version))
}
