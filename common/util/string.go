package util

import (
	"fmt"
	"strconv"
)

func StrLen(str string) int {
	return len([]rune(str))
}

func StringNormalize(str string, min int, max int) string {
	fmtStr := fmt.Sprintf("%%-%d.%ds", min, max)
	return fmt.Sprintf(fmtStr, str)
}

func Atoi(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

func GetOrDefault(s string, def string) string {
	if s == "" {
		return def
	}
	return s
}

func LevenshteinDistance(s1 string, s2 string) int {
	// support unicode
	r1 := []rune(s1)
	r2 := []rune(s2)
	r1l := len(r1)
	r2l := len(r2)
	if r1l == 0 || r2l == 0 {
		return Max(r1l, r2l)
	}
	previous := make([]int, r2l+1)
	current := make([]int, r2l+1)

	for i := 0; i <= r2l; i++ {
		previous[i] = i
	}

	for i := 1; i <= r1l; i++ {
		current[0] = i
		for j := 1; j <= r2l; j++ {
			subCost := 1
			if r1[i-1] == r2[j-1] {
				subCost = 0
			}
			// current[j] = min( insertCost,deleteCost, subCost)
			current[j] = Min(current[j-1]+1, previous[j]+1, previous[j-1]+subCost)
		}
		current, previous = previous, current
	}
	return previous[r2l]
}

func WeightedLevenshteinDistance(s1 string, s2 string, ins, del, repl int) int {
	// support unicode
	r1 := []rune(s1)
	r2 := []rune(s2)
	r1l := len(r1)
	r2l := len(r2)
	if r1l == 0 || r2l == 0 {
		return Max(r1l, r2l)
	}
	previous := make([]int, r2l+1)
	current := make([]int, r2l+1)

	for i := 0; i <= r2l; i++ {
		previous[i] = i
	}

	for i := 1; i <= r1l; i++ {
		current[0] = i
		for j := 1; j <= r2l; j++ {
			subCost := 1
			if r1[i-1] == r2[j-1] {
				subCost = 0
			}
			// current[j] = min( insertCost,deleteCost, subCost)
			current[j] = Min(current[j-1]+1*ins, previous[j]+1*del, previous[j-1]+subCost*repl)
		}
		current, previous = previous, current
	}
	return previous[r2l]
}

func LongestCommonString(s1 string, s2 string) string {
	// support unicode
	r1 := []rune(s1)
	r2 := []rune(s2)
	r1l := len(r1)
	r2l := len(r2)
	if r1l == 0 || r2l == 0 {
		return ""
	}
	previous := make([]int, r2l+1)
	current := make([]int, r2l+1)
	max := 0
	maxIndex := 0
	for i := 1; i <= r1l; i++ {
		for j := 1; j <= r2l; j++ {
			if r1[i-1] == r2[j-1] {
				current[j] = previous[j-1] + 1
				if current[j] > max {
					max = current[j]
					maxIndex = i
				}
			} else {
				current[j] = 0
			}
		}
		current, previous = previous, current
	}
	return string(r1[maxIndex-max : maxIndex])
}
