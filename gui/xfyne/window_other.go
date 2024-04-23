//go:build !darwin && !windows && !linux
// +build !darwin,!windows,!linux

package xfyne

import "fyne.io/fyne/v2"

func GetWindowHandle(window fyne.Window) uintptr {
	return 0
}
