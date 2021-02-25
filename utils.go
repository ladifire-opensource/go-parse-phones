package goparsephone

import (
	"reflect"
	"strconv"
)

func ArrayFilter(input []int, callback func(int) bool) []int {
	var val = reflect.ValueOf(input)
	var res []int

	for i := 0; i < val.Len(); i++ {
		v := val.Index(i).Int()

		if callback(int(v)) {
			res = append(res, int(v))
		}
	}

	return res
}

func IsNumeric(s string) bool {
	result := false

	if _, err := strconv.Atoi(s); err == nil {
		result = true
	}

	return result
}
