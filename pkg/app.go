package pkg

import (
	"log"
	"os"
	"strings"

	"github.com/chzyer/readline"
)

type CSVQL struct {
	HistoryFile  string
	DatabasePath string
	Sql          string
	Tokens       []Token
	Ast          AST
	Result       [][]string
	Duration     float64
	Readline     *readline.Instance
	Variables    map[string]string
	Error        error
}

func NewQuery(sql string, databasePath string) CSVQL {
	_, err := os.ReadDir(databasePath)

	if err != nil {
		log.Fatalln("Directory not found. Please try again.")
	}
	replacer := strings.NewReplacer("\n", " ", "\t", " ")
	sql = replacer.Replace(sql)
	query := CSVQL{
		Sql:          sql,
		DatabasePath: databasePath,
		HistoryFile:  "/tmp/readline-multiline",
		Tokens:       []Token{},
		Ast:          AST{},
		Result:       [][]string{},
		Variables:    map[string]string{},
	}

	return query
}

func (ql *CSVQL) Execute() {
	ql.Tokenizer()
	ql.BuildAST()
	ql.ExecuteAST()
}
