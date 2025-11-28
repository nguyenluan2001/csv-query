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

func SelectField(row []string, headerIndex map[string]int, columns []Token) []string {
	newRow := []string{}
	if columns[0].Type == TokenStar {
		return row
	}
	for _, col := range columns {
		headerIdx := headerIndex[fmt.Sprintf("%v", col.Value)]
		newRow = append(newRow, row[headerIdx])
	}
	return newRow
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
	header := []string{}
	headerIndex := map[string]int{}
	result := [][]string{}
	for {
		row, err := csvReader.Read()
		if err == io.EOF {
			break // End of file reached
		}
		if err != nil {
			// Handle error during reading
		}

		if len(header) == 0 {
			header = row
			headerIndex = ParseHeaderIndex(header)
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

		//Handle SELECT statement
		result = append(result, SelectField(row, headerIndex, ast.Columns))
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
	fmt.Println("result:", result)
	return nil

}
