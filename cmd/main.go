package main

import (
	"fmt"

	"github.com/nguyenluan2001/csv-query/utils"
)

func main() {
	// sql := "SELECT name,salary, age FROM employees WHERE salary > 50000 AND age >=30"
	sql := "SELECT name,salary, age FROM employees WHERE salary > 50000"
	tokens := utils.Tokenizer(sql)
	fmt.Println("tokens", tokens)

	ast, err := utils.BuildAST(tokens)
	fmt.Println("ast", ast.WhereCombine.Left)
	fmt.Println("ast", ast.WhereCombine.Right)
	fmt.Println("err", err)

	utils.Execute(ast)
}
