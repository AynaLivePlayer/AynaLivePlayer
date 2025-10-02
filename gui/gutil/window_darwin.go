//go:build darwin

package gutil

import (
	"fyne.io/fyne/v2"
)

func GetWindowHandle(window fyne.Window) uintptr {
	glfwWindow := getGlfwWindow(window)
	if glfwWindow == nil {
		return 0
	}
	return uintptr(glfwWindow.GetCocoaWindow())
}
