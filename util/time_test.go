package util

import (
	"fmt"
	"testing"
)

func TestFormatTime(t *testing.T) {
	fmt.Println(FormatTime(60 * 60))
}
