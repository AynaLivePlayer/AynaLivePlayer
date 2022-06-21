package util

import "fmt"

func FormatTime(time int) string {
	return fmt.Sprintf("%01d:%02d", time/60, time%60)
}
