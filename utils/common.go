package utils

import (
	"fmt"
	"strconv"
)

func StringToInt(s interface{}) (int, error) {
	return strconv.Atoi(fmt.Sprintf("%v", s))
}

func BooleanToInt(value bool) int {
	if value {
		return 1
	}
	return 0
}

func Stringify(value interface{}) string {
	return fmt.Sprintf("%v", value)
}
