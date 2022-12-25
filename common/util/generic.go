package util

import (
	"golang.org/x/exp/constraints"
)

func Slice[T any](arr []T, from int, to int) []T {
	l := len(arr)
	to = Min(to, l)
	from = Min(from, l)
	if to <= 0 {
		to = l + to
		if to <= 0 {
			return []T{}
		}
	}
	if from < 0 {
		from = l + from
		if from < 0 {
			from = 0
		}
	}
	if to <= from {
		return []T{}
	}
	return arr[from:to]
}

func SliceCopy[T any](src []T) []T {
	x := make([]T, len(src))
	copy(x, src)
	return x
}

func SliceContains[T comparable](s []T, e T) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func TernaryOp[T any](cond bool, a T, b T) T {
	if cond {
		return a
	}
	return b
}

func Min[T constraints.Ordered](arr ...T) T {
	min := arr[0]
	for _, a := range arr {
		if a < min {
			min = a
		}
	}
	return min
}

func Max[T constraints.Ordered](arr ...T) T {
	max := arr[0]
	for _, a := range arr {
		if a > max {
			max = a
		}
	}
	return max
}
