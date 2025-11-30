package utils

import (
	"fmt"
	"slices"
)

func Sum(collector []string) string {
	sum := 0
	for _, val := range collector {
		valNumber, _ := StringToInt(Stringify(val))
		sum += valNumber
	}
	return Stringify(sum)
}

func Count(collector []string) string {
	return Stringify(len(collector))
}

func Average(collector []string) string {
	sum := 0
	for _, val := range collector {
		valNumber, _ := StringToInt(Stringify(val))
		sum += valNumber
	}
	return Stringify(sum / len(collector))
}

func Max(collector []string) string {
	slices.SortFunc(collector, func(val1, val2 string) int {
		return CompareProxy(val2, val1)
	})
	if len(collector) == 0 {
		return ""
	}
	return Stringify(collector[0])
}

func Min(collector []string) string {
	slices.SortFunc(collector, func(val1, val2 string) int {
		return CompareProxy(val1, val2)
	})
	if len(collector) == 0 {
		return ""
	}
	return Stringify(collector[0])
}

func GetAggregateFnName(token Token) string {
	switch token.Type {
	case TokenSum:
		{
			return "SUM"
		}
	case TokenCount:
		{
			return "COUNT"
		}
	case TokenMax:
		{
			return "MAX"
		}
	case TokenMin:
		{
			return "MIN"
		}
	case TokenAverage:
		{
			return "AVG"
		}
	default:
		{
			return ""
		}
	}
}

func GenerateAggregateColumnName(token Token) string {
	return fmt.Sprintf("%v_%v", GetAggregateFnName(token), token.Value)
}
