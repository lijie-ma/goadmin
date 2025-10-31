package util

import (
	"cmp"
	"slices"
)

// Unique 返回去重后的新切片（仅保留一个相同值的副本）
// 内部使用 slices.Sort + slices.Compact
func Unique[T cmp.Ordered](s []T) []T {
	if len(s) == 0 {
		return s
	}
	slices.Sort(s)
	return slices.Compact(s)
}

// Difference 返回存在于 a 中但不在 b 中的元素集合
func Difference[T comparable](a, b []T) []T {
	if len(a) == 0 {
		return nil
	}
	if len(b) == 0 {
		return slices.Clone(a)
	}

	bset := make(map[T]struct{}, len(b))
	for _, v := range b {
		bset[v] = struct{}{}
	}

	var diff []T
	for _, v := range a {
		if _, found := bset[v]; !found {
			diff = append(diff, v)
		}
	}
	return diff
}
