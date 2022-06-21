package util

import (
	"fmt"
	"strconv"
)

func SliceString(str string, from int, to int) (string, bool) {
	sList := []rune(str)
	if to <= 0 {
		to = len(sList) + to
	}
	if from >= len(sList) || to > len(sList) {
		return "", false
	}
	return string(sList[from:to]), true
}

func LenString(str string) int {
	return len([]rune(str))
}

func StringNormalize(str string, min int, max int) string {
	fmtStr := fmt.Sprintf("%%-%d.%ds", min, max)
	return fmt.Sprintf(fmtStr, str)
}

func StringSliceContains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func StringToInt(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}
