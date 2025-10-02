//go:build !darwin && !windows && !linux

package gutil

import "fyne.io/fyne/v2"

func GetWindowHandle(window fyne.Window) uintptr {
	return 0
}
