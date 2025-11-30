package main

import (
	"fmt"

	"github.com/kr/pretty"
	"github.com/nguyenluan2001/csv-query/utils"
)

func main() {
	// sql := "SELECT name,salary, age FROM employees WHERE salary > 50000 AND age >=30"
	// sql := "SELECT * FROM employees WHERE salary >= 50000 AND age >= 30 AND id >= 2"
	// sql := "SELECT * FROM employees WHERE age <> 40"
	// sql := "SELECT * FROM employees WHERE age BETWEEN 30 AND 40 AND salary > 50000"
	// sql := "SELECT * FROM employees WHERE age BETWEEN 30 AND 40"
	// sql := "SELECT * FROM employees WHERE salary BETWEEN 50000 AND 80000"
	// sql := "SELECT * FROM employees WHERE salary BETWEEN 50000 AND 80000 AND age BETWEEN 35 AND 40"
	// sql := "SELECT * FROM employees WHERE salary BETWEEN 51000 AND 80000 AND age < 30"
	// sql := "SELECT * FROM employees WHERE name = 'Eve' "
	// sql := "SELECT * FROM employees WHERE id IN (1, 2) AND salary > 50000"
	// sql := "SELECT * FROM employees WHERE name IN ('Alice', 'Dan', 'Eve') OR name='Frank' ORDER BY salary DESC"
	// sql := "SELECT * FROM employees WHERE dept IN ('eng', 'ops') AND age BETWEEN 30 AND 40"
	// sql := "SELECT * FROM employees WHERE age > 30 ORDER BY salary"
	// sql := "SELECT * FROM employees ORDER BY id DESC"
	// sql := "SELECT * FROM employees ORDER BY salary DESC, age"
	// sql := "SELECT * FROM employees ORDER BY name DESC, age ASC"
	// sql := "SELECT * FROM employees LIMIT 1"
	sql := "SELECT AVG(salary), COUNT(id), AVG(age), dept FROM employees GROUP BY dept"
	tokens := utils.Tokenizer(sql)

	ast, err := utils.BuildAST(tokens)
	fmt.Println("tokens", tokens)
	fmt.Printf("ast:%# v", pretty.Formatter(ast))
	// fmt.Println("ast", ast.WhereCombine.Right)
	fmt.Println("err", err)

	utils.Execute(ast)
}
