//go:build linux
// +build linux

package xfyne

import (
	"fyne.io/fyne/v2"
)

func GetWindowHandle(window fyne.Window) uintptr {
	glfwWindow := getGlfwWindow(window)
	if glfwWindow == nil {
		return 0
	}
	return uintptr(glfwWindow.GetX11Window())
}
