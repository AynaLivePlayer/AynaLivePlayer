package xfyne

import (
	"fyne.io/fyne/v2"
	"github.com/go-gl/glfw/v3.3/glfw"
	"reflect"
	"unsafe"
)

// getGlfwWindow returns the glfw.Window pointer from a fyne.Window.
// very unsafe and ugly hacks. but it works.
func getGlfwWindow(window fyne.Window) *glfw.Window {
	rv := reflect.ValueOf(window)
	if rv.Type().String() != "*glfw.window" {
		return nil
	}
	rv = rv.Elem()
	var glfwWindowPtr uintptr = rv.FieldByName("viewport").Pointer()
	for glfwWindowPtr == 0 {
		glfwWindowPtr = rv.FieldByName("viewport").Pointer()
	}
	return (*glfw.Window)(unsafe.Pointer(glfwWindowPtr))
}
