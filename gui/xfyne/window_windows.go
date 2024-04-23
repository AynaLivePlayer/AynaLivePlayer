//go:build windows
// +build windows

package xfyne

import (
	"fyne.io/fyne/v2"
	"unsafe"
)

func GetWindowHandle(window fyne.Window) uintptr {
	glfwWindow := getGlfwWindow(window)
	if glfwWindow == nil {
		return 0
	}
	return uintptr(unsafe.Pointer(glfwWindow.GetWin32Window()))
}
