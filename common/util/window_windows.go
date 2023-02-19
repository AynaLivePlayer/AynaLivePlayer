package util

import (
	"golang.org/x/sys/windows"
	"strings"
	"syscall"
	"unsafe"
)

var user32dll = windows.MustLoadDLL("user32.dll")

func getWindowHandle(title string) uintptr {
	var the_handle uintptr
	window_byte_name := []byte(title)

	// Windows will loop over this function for each window.
	wndenumproc_function := syscall.NewCallback(func(hwnd uintptr, lparam uintptr) uintptr {
		// Allocate 100 characters so that it has something to write to.
		var filename_data [100]uint16
		user32dll.MustFindProc("GetWindowTextW").Call(hwnd, uintptr(unsafe.Pointer(&filename_data)), uintptr(100))

		// If there's a match, save the value and return 0 to stop the iteration.
		if strings.Contains(windows.UTF16ToString(filename_data[:]), string(window_byte_name)) {
			the_handle = hwnd
			return 0
		}
		return 1
	})
	// Call the above looping function.
	user32dll.MustFindProc("EnumWindows").Call(wndenumproc_function, uintptr(0))

	return the_handle
}
