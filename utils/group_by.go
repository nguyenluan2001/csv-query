package utils

import (
	"slices"
	"strings"
)

// GROUP BY utils
func ProcessGroupByPerRow(row []string, headerRow []string, headerIndex map[string]int, tokens []Token) ([]string, []string, []string, []string, string) {
	keyArr := []string{"groupBy"}
	fieldIndexes := []int{}
	groupByFields := []string{}
	groupByData := []string{}
	otherFields := []string{}
	otherData := []string{}
	for _, token := range tokens {
		headerIdx := headerIndex[Stringify(token.Value)]
		fieldIndexes = append(fieldIndexes, headerIdx)
	}
	for i, val := range row {
		if !slices.Contains(fieldIndexes, i) {
			otherData = append(otherData, val)
			otherFields = append(otherFields, headerRow[i])
			continue
		}
		keyArr = append(keyArr, val)
		groupByFields = append(groupByFields, headerRow[i])
		groupByData = append(groupByData, val)
	}
	return groupByFields, groupByData, otherFields, otherData, strings.Join(keyArr, "_")
}

func AppendGroupByData(groupByMap map[string]GroupByData, key string, otherFields []string, otherData []string) bool {
	_, ok := groupByMap[key]
	if !ok {
		groupByMap[key] = map[string]Collector{}
	}
	isGrouped := false
	for i, field := range otherFields {
		_, ok := groupByMap[key][field]
		if !ok {
			groupByMap[key][field] = Collector{otherData[i]}
		} else {
			groupByMap[key][field] = append(groupByMap[key][field], otherData[i])
			isGrouped = true
		}
	}
	return isGrouped
}

func BuildGroupByRow(groupByFields []string, groupByData []string, headerRow []string, key string) []string {
	row := []string{}
	groupByPointer := 0
	for _, val := range headerRow {
		if !slices.Contains(groupByFields, val) {
			row = append(row, key)
		} else {
			row = append(row, groupByData[groupByPointer])
			groupByPointer++
		}
	}
	return row
}

func GetAggregateCollector(groupByMap map[string]GroupByData, groupByKey, field string) []string {
	collector, ok := groupByMap[groupByKey][field]
	if !ok {
		return []string{}
	}
	return collector
}

func GetAggregateValue(token Column, collector []string) string {
	switch token.Type {
	case TokenSum:
		{
			return Sum(collector)
		}
	case TokenCount:
		{
			return Count(collector)
		}
	case TokenAverage:
		{
			return Average(collector)
		}
	case TokenMax:
		{
			return Max(collector)
		}
	case TokenMin:
		{
			return Min(collector)
		}
	default:
		{
			panic("Syntax error")
		}
	}
}
