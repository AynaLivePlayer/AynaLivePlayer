package logger

import "fmt"

const (
	LogColorBlack Color = iota + 30
	LogColorRed
	LogColorGreen
	LogColorYellow
	LogColorBlue
	LogColorMagenta
	LogColorCyan
	LogColorWhite
)

// Color represents a text color.
type Color uint8

// Add adds the coloring to the given string.
func (c Color) Add(s string) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", uint8(c), s)
}

var LogColorMap = map[LogLevel]Color{
	LogLevelError: LogColorRed,
	LogLevelWarn:  LogColorYellow,
	LogLevelInfo:  LogColorCyan,
	LogLevelDebug: LogColorWhite,
}
