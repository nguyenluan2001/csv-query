package main

import (
	"fmt"
	"os"
	"path"

	"github.com/nguyenluan2001/csv-query/utils"
)

func main() {
	// sql := "SELECT name AS name__123, salary, age FROM employees WHERE salary > 50000 AND age >=30"
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
	// sql := "SELECT AVG(salary) AS avg_salary123, COUNT(id), AVG(age) AS avg_age, dept FROM employees GROUP BY dept"
	// sql := "SELECT AVG(salary) AS avg_salary, COUNT(id) AS count_id, AVG(age) FROM employees GROUP BY dept"
	// sql := "SELECT COUNT(employee_id), department_id FROM employees_2 GROUP BY department_id"
	sql := "SELECT employees_2.email, employee_id, email, deparment_name FROM employees_2 JOIN department ON employees_2.department_id = department.department_id"
	// sql := "employees_2 LEFT JOIN department ON employees_2.department_id = department.department_id AND employees_2.first_name = 'John' "

	// JoinExpr{
	// 	Type: "LEFT",
	// 	Left: JoinExpr{
	// 		Tye:   "INNER",
	// 		Left:  "employees",
	// 		Right: "departments",
	// 		Condition: OnExpr{
	// 			Left:  "employees.department_id",
	// 			Op:    "=",
	// 			Right: "departments.id",
	// 		},
	// 	},
	// 	Right: "job",
	// 	Condition: OnExpr{
	// 		Left:  "employees.job_id",
	// 		Op:    "=",
	// 		Right: "jobs.id",
	// 	},
	// }

	tokens := utils.Tokenizer(sql)
	fmt.Println("tokens", tokens)

	// p := utils.NewParserFrom(tokens)

	// (*p).RegisterPrefix(utils.TokenNumber, 0)
	// (*p).RegisterPrefix(utils.TokenString, 0)
	// (*p).RegisterPrefix(utils.TokenIdent, 0)
	// (*p).RegisterInfix(utils.TokenOr, 100)
	// (*p).RegisterInfix(utils.TokenAnd, 200)
	// (*p).RegisterInfix(utils.TokenEqual, 500)
	// left, _ := (*p).ParserFromExpression(0)
	// fmt.Printf("left:%# v", pretty.Formatter(left))

	ast, err := utils.BuildAST(tokens)
	// fmt.Printf("ast:%# v", pretty.Formatter(ast))
	utils.PrintPretty("AST", ast)
	fmt.Println("err", err)

	cwd, _ := os.Getwd()
	databasePath := path.Join(cwd, "../database")
	utils.Execute(ast, databasePath)
}
