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

// GROUP BY
type Collector []string
type GroupByData map[string]Collector
type GroupByRow struct {
	Data  map[string][][]string
	Order map[string][]int
}

// JOIN
type ScanTableInfo struct {
	HeaderRow   []string
	HeaderIndex map[string]int
	Rows        [][]string
	GroupByRow  GroupByRow
}

type JoinHeaderIndex map[string][]int

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

// Select field in normal mode
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

// Select field in GROUP BY mode
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

// Select field in JOIN mode
func SelectJoinField(row []string, headerIndex map[string]JoinHeaderIndex, columns []Column) ([]string, []string) {
	header := []string{}
	newRow := []string{}
	for _, col := range columns {
		valueArr, isArray := col.Value.([]string)

		colName := Stringify(col.Value)
		var headerIdx []int
		if !isArray {
			headerIdx = headerIndex["global"][colName]
			if len(headerIdx) > 1 {
				panic(fmt.Sprintf(`Column reference "%v" is ambiguous `, colName))
			}
		} else {
			headerIdx = headerIndex[valueArr[0]][valueArr[1]]
			colName = valueArr[1]
		}
		header = append(header, colName)
		newRow = append(newRow, row[headerIdx[0]])
	}
	// return []string{}, []string{}
	return newRow, header
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

func ScanTable(databasePath, name string, groupBy string) ScanTableInfo {
	filepath := path.Join(databasePath, fmt.Sprintf("%v.csv", name))
	scanTableInfo := ScanTableInfo{}
	groupByRow := GroupByRow{
		Data:  map[string][][]string{},
		Order: make(map[string][]int),
	}

	file, err := os.Open(filepath)

	if err != nil {
		panic(fmt.Sprintf("Table not found: %s", filepath))
	}
	defer file.Close()

	csvReader := csv.NewReader(file)
	rowCount := 0
	for {
		row, err := csvReader.Read()
		if err == io.EOF {
			break // End of file reached
		}
		if err != nil {
			// Handle error during reading
		}

		if len(scanTableInfo.HeaderRow) == 0 {
			// originalHeaderRow = row
			scanTableInfo.HeaderRow = row
			scanTableInfo.HeaderIndex = ParseHeaderIndex(scanTableInfo.HeaderRow)
			continue
		}

		if len(groupBy) > 0 {
			headerIdx, _ := scanTableInfo.HeaderIndex[groupBy]
			groupByValue := row[headerIdx]

			_, ok := groupByRow.Data[groupByValue]

			if !ok {
				groupByRow.Data[groupByValue] = [][]string{row}
				groupByRow.Order[groupByValue] = []int{rowCount}
			} else {
				groupByRow.Data[groupByValue] = append(groupByRow.Data[groupByValue], row)
				groupByRow.Order[groupByValue] = append(groupByRow.Order[groupByValue], rowCount)
			}
		}

		// result = append(result, row)
		scanTableInfo.Rows = append(scanTableInfo.Rows, row)
		rowCount++
	}
	scanTableInfo.GroupByRow = groupByRow
	return scanTableInfo
}

func HandleJoinTable(databasePath string, joinExpr JoinExpr) ([][]string, map[string]JoinHeaderIndex, []string) {
	leftTableToken, isLeftSingleTable := joinExpr.Left.(Token)
	rightTableToken, isRightSingleTable := joinExpr.Right.(Token)
	leftScanTableInfo := ScanTableInfo{}
	rightScanTableInfo := ScanTableInfo{}
	joinHeader := []string{}
	joinHeaderIndex := map[string]JoinHeaderIndex{
		"global": make(JoinHeaderIndex),
	}
	columnIdx := 0

	condition, _ := joinExpr.Condition.(BinaryExpr)

	if isLeftSingleTable {
		leftCondition, _ := condition.Left.(TableIdentifier)
		tableName := Stringify(leftTableToken.Value)
		groupByField := Stringify(leftCondition.Field.Value)

		leftScanTableInfo = ScanTable(databasePath, tableName, groupByField)

		joinHeader = append(joinHeader, leftScanTableInfo.HeaderRow...)
		joinHeaderIndex[tableName] = make(JoinHeaderIndex)
		for key, value := range leftScanTableInfo.HeaderIndex {
			joinHeaderIndex["global"][key] = []int{value}
			joinHeaderIndex[tableName][key] = []int{value}
			columnIdx++
		}
		PrintPretty("leftScanTableInfo", leftScanTableInfo)
	} else {
		joinExpr, _ := joinExpr.Left.(JoinExpr)
		_result, _joinHeaderIndex, _joinHeader := HandleJoinTable(databasePath, joinExpr)
		joinHeader = append(joinHeader, _joinHeader...)
		leftScanTableInfo.Rows = _result
		leftScanTableInfo.HeaderRow = _joinHeader
		PrintPretty("_joinHeaderIndex", _joinHeaderIndex)
		for groupKey, groupHash := range _joinHeaderIndex {
			for key, value := range groupHash {
				joinHeaderIndex[groupKey][key] = value
				joinHeaderIndex[groupKey][key] = value
				columnIdx++
			}
		}

	}

	if isRightSingleTable {
		rightCondition, _ := condition.Right.(TableIdentifier)
		tableName := Stringify(rightTableToken.Value)
		groupByField := Stringify(rightCondition.Field.Value)

		rightScanTableInfo = ScanTable(databasePath, tableName, groupByField)

		joinHeader = append(joinHeader, rightScanTableInfo.HeaderRow...)
		joinHeaderIndex[tableName] = make(JoinHeaderIndex)
		for key, value := range rightScanTableInfo.HeaderIndex {
			_, ok := joinHeaderIndex["global"][key]
			if !ok {
				joinHeaderIndex["global"][key] = []int{value + columnIdx}
			} else {
				joinHeaderIndex["global"][key] = append(joinHeaderIndex["global"][key], value+columnIdx)
			}
			joinHeaderIndex[tableName][key] = []int{value + columnIdx}
		}
		PrintPretty("rightScanTableInfo", rightScanTableInfo)
	}

	// result := [][]string{joinHeader}
	result := make([][]string, len(leftScanTableInfo.Rows))
	for key, rows := range leftScanTableInfo.GroupByRow.Data {
		rightRows := rightScanTableInfo.GroupByRow.Data[key]
		rowOrder := leftScanTableInfo.GroupByRow.Order[key]
		for i, lRow := range rows {
			for _, rRow := range rightRows {
				fullRow := append([]string{}, lRow...)
				fullRow = append(fullRow, rRow...)
				// result = append(result, fullRow)
				result[rowOrder[i]] = fullRow
			}
		}
	}
	result = append([][]string{joinHeader}, result...)
	return result, joinHeaderIndex, joinHeader
}

func Execute(ast AST, databasePath string) error {
	originalHeaderRow := []string{}
	headerRow := []string{}
	headerIndex := map[string]int{}
	result := [][]string{}
	groupByMap := map[string]GroupByData{}

	groupByTable := [][]string{}
	groupByHeader := []string{}

	joinTable := [][]string{}
	joinHeaderIndex := map[string]JoinHeaderIndex{}
	joinHeader := []string{}

	joinExpr, isJoin := ast.From.(JoinExpr)
	if !isJoin {
		from, _ := ast.From.(Token)
		// Handle FROM statement
		filepath := path.Join(databasePath, fmt.Sprintf("%v.csv", from.Value))
		fmt.Println(filepath)
		file, err := os.Open(filepath)
		if err != nil {
			log.Fatalln("File not found.")
		}

		csvReader := csv.NewReader(file)
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
	} else {
		joinTable, joinHeaderIndex, _ = HandleJoinTable(databasePath, joinExpr)
		PrintPretty("joinHeaderIndex", joinHeaderIndex)
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
	} else if isJoin {
		PrintPretty("joinTable", joinTable)
		for i := 1; i < len(joinTable); i++ {
			row := joinTable[i]
			newRow, newHeader := SelectJoinField(row, joinHeaderIndex, ast.Columns)
			result = append(result, newRow)

			if len(joinHeader) == 0 {
				joinHeader = newHeader
				headerIndex = ParseHeaderIndex(newHeader)
			}

		}
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
	} else if isJoin {
		result = append([][]string{joinHeader}, result...)
	} else {
		result = append([][]string{headerRow}, result...)
	}
	PrintPretty("GroupByMap", groupByMap)
	PrintPretty("GroupByTable", groupByTable)
	PrintPretty("result: ", result)
	// fmt.Println("result:", result)
	return nil

}
