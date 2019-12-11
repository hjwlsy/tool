package tool

import (
	"strings"
)

func String2IntArray(ids string) []int {
	if ids == "" {
		return make([]int, 0)
	}
	sl := strings.Split(ids, ",")
	arr := make([]int, len(sl))
	for k, v := range sl {
		arr[k] = StringToInt(v)
	}
	return arr
}

func InArrayInt(n int, h []int) bool {
	for _, v := range h {
		if v == n {
			return true
		}
	}
	return false
}
