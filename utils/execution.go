package utils

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"path"
	"slices"
	"strconv"
	"strings"
)

type Collector []string
type GroupByData map[string]Collector

func ParseHeaderIndex(header []string) map[string]int {
	headerIndex := map[string]int{}
	for i, val := range header {
		headerIndex[val] = i
	}
	return headerIndex
}

func Compare(a, b int, op string) bool {
	switch op {
	case ">":
		{
			return a > b
		}
	case ">=":
		{
			return a >= b
		}
	case "<":
		{
			return a < b
		}
	case "<=":
		{
			return a <= b
		}
	case "!=":
		{
			return a != b
		}
	case "=":
		{
			return a == b
		}
	default:
		{
			return false
		}
	}
}

func FilterRow(row []string, headerIndex map[string]int, condition BinaryExpr) bool {
	// headerIdx := headerIndex[fmt.Sprintf("%v", condition.Left.Value)]
	// leftVal, _ := strconv.Atoi(fmt.Sprintf("%v", row[headerIdx]))
	// rightVal, _ := strconv.Atoi(fmt.Sprintf("%v", condition.Right.Value))
	// opVal := fmt.Sprintf("%v", condition.Op.Value)
	// return Compare(int(leftVal), int(rightVal), opVal)
	return true
}

func SelectField(row []string, headerIndex map[string]int, columns []Column, headerRow []string) ([]string, []string) {
	newRow := []string{}
	newHeader := []string{}
	if columns[0].Type == TokenStar {
		return row, headerRow
	}
	for _, col := range columns {
		headerIdx := headerIndex[fmt.Sprintf("%v", col.Value)]
		newRow = append(newRow, row[headerIdx])
		newHeader = append(newHeader, Stringify(col.Value))
		if len(col.Alias) > 0 {
			newHeader[len(newHeader)-1] = Stringify(col.Alias)
		}
	}
	return newRow, newHeader
}

func SelectGroupByField(row []string, headerIndex map[string]int, columns []Column, groupByMap map[string]GroupByData) ([]string, []string) {
	newRow := []string{}
	newHeader := []string{}
	for _, col := range columns {
		headerIdx := headerIndex[fmt.Sprintf("%v", col.Value)]
		if !IsAggregateFn(col) {
			newHeader = append(newHeader, Stringify(col.Value))
			newRow = append(newRow, row[headerIdx])
		} else {
			// aggregateFn := GetAggregateFn(col)
			collector := GetAggregateCollector(groupByMap, row[headerIdx], Stringify(col.Value))
			PrintPretty("collector", collector)
			aggregateVal := GetAggregateValue(col, collector)
			newRow = append(newRow, aggregateVal)
			newHeader = append(newHeader, GenerateAggregateColumnName(col))
		}

		if len(col.Alias) > 0 {
			newHeader[len(newHeader)-1] = Stringify(col.Alias)
		}
	}
	return newRow, newHeader
}

// row1[field] - row2[field] < 0: asc
// row2[field] - row1[field] > 0: desc
// row1[field] - row2[field] == 0: remain
// Process condition one by one. If row1[field]==row2[field] => process next condition
func CompareNumber(a, b int) int {
	return a - b
}

func CompareString(a, b string) int {
	return strings.Compare(a, b)
}

func CompareProxy(a, b interface{}) int {
	aNumber, aErr := strconv.Atoi(Stringify(a))
	bNumber, bErr := strconv.Atoi(Stringify(b))

	if aErr != nil || bErr != nil {
		return CompareString(Stringify(a), Stringify(b))
	}
	return CompareNumber(aNumber, bNumber)
}

func OrderByComparator(conditions []OrderBySingle, headerIndex map[string]int, row1, row2 []string) int {
	for _, condition := range conditions {
		field := condition.Field
		direction := condition.Direction
		fieldIdx := headerIndex[field]
		row1Val := Stringify(row1[fieldIdx])
		row2Val := Stringify(row2[fieldIdx])

		if row1Val == row2Val {
			continue
		}

		if direction == TokenAsc {
			return CompareProxy(row1Val, row2Val)
		}

		return CompareProxy(row2Val, row1Val)
	}
	return 0
}

func Execute(ast AST) error {

	// Handle FROM statement
	cwd, _ := os.Getwd()
	filepath := path.Join(cwd, "..", fmt.Sprintf("%v.csv", ast.From.Value))
	fmt.Println(filepath)
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatalln("File not found.")
	}

	csvReader := csv.NewReader(file)
	originalHeaderRow := []string{}
	headerRow := []string{}
	headerIndex := map[string]int{}
	result := [][]string{}
	groupByMap := map[string]GroupByData{}
	groupByTable := [][]string{}
	groupByHeader := []string{}
	for {
		row, err := csvReader.Read()
		if err == io.EOF {
			break // End of file reached
		}
		if err != nil {
			// Handle error during reading
		}

		if len(originalHeaderRow) == 0 {
			originalHeaderRow = row
			headerIndex = ParseHeaderIndex(originalHeaderRow)
			continue
		}

		var isValid interface{}

		//Handle WHERE statement
		if ast.Where != nil {
			isValid = Eval(ast.Where, row, headerIndex)
		}
		if isValid == 0 {
			continue
		}

		if ast.GroupBy == nil { // In case no GROUP BY => Process SELECT columns
			//Handle SELECT statement
			newRow, newHeader := SelectField(row, headerIndex, ast.Columns, originalHeaderRow)
			result = append(result, newRow)

			headerRow = newHeader
		} else { // In case GROUP BY => Process groupByMap and groupByTable
			groupByFields, groupByData, otherFields, otherData, key := ProcessGroupByPerRow(row, originalHeaderRow, headerIndex, ast.GroupBy)
			isGrouped := AppendGroupByData(groupByMap, key, otherFields, otherData)
			if !isGrouped {
				groupByTable = append(groupByTable, BuildGroupByRow(groupByFields, groupByData, originalHeaderRow, key))
			}
			fmt.Println("key", groupByFields, groupByData, otherFields, otherData, key)
		}
	}

	if ast.GroupBy != nil {
		for _, row := range groupByTable {
			newRow, newHeader := SelectGroupByField(row, headerIndex, ast.Columns, groupByMap)
			result = append(result, newRow)

			if len(groupByHeader) == 0 {
				groupByHeader = newHeader
			}
		}
		// result = append([][]string{groupByHeader}, result...)
	}

	if ast.OrderBy != nil {
		slices.SortFunc(result, func(row1, row2 []string) int {
			return OrderByComparator(ast.OrderBy, headerIndex, row1, row2)
		})
	}

	if ast.Limit > 0 {
		maxLimit := int(math.Min(float64(ast.Limit), float64(len(result))))
		result = result[0:maxLimit]
	}

	//Adding result header
	if ast.GroupBy != nil {
		result = append([][]string{groupByHeader}, result...)
	} else {
		result = append([][]string{headerRow}, result...)
	}
	PrintPretty("GroupByMap", groupByMap)
	PrintPretty("GroupByTable", groupByTable)
	PrintPretty("result: ", result)
	// fmt.Println("result:", result)
	return nil

}
