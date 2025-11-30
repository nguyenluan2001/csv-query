package utils

import (
	"fmt"
	"strconv"

	"github.com/kr/pretty"
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

func PrintPretty(key string, value interface{}) {
	fmt.Printf("%v: %# v \n", key, pretty.Formatter(value))
}
