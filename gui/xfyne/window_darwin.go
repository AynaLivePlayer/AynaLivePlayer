//go:build darwin
// +build darwin

package xfyne

func GetWindowHandle(window fyne.Window) uintptr {
	glfwWindow := getGlfwWindow(window)
	if glfwWindow == nil {
		return 0
	}
	return uintptr(glfwWindow.GetCocoaWindow())
}
