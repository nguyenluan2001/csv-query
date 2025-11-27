package utils

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"path"
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

func Execute(ast AST) error {
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
		}

		// isValid := FilterRow(row, headerIndex, ast.Where)
		isValid := Eval(ast.Where, row, headerIndex)
		if isValid == 0 {
			continue
		}
		result = append(result, SelectField(row, headerIndex, ast.Columns))
	}
	fmt.Println("result:", result)
	return nil

}
