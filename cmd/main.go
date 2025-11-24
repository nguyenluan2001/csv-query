package main

import (
	"fmt"

	"github.com/nguyenluan2001/csv-query/utils"
)

func main() {
	sql := "SELECT id, salary FROM employees WHERE salary > 50000"
	tokens := utils.Tokenizer(sql)
	fmt.Println("tokens", tokens)
}
