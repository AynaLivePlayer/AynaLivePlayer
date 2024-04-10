package util

import "fmt"

// FormatTime formats time in seconds to string in format "m:ss"
func FormatTime(sec int) string {
	return fmt.Sprintf("%01d:%02d", sec/60, sec%60)
}
