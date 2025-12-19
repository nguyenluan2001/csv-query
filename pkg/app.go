package pkg

import (
	"encoding/csv"
	"io"
	"log"
	"os"

	"github.com/chzyer/readline"
	"github.com/nguyenluan2001/csv-query/utils"
)

type CSVQL struct {
	HistoryFile  string
	VariableFile string
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
	variableFile := "/tmp/csvql_variable"

	query := CSVQL{
		Sql:          sql,
		DatabasePath: databasePath,
		HistoryFile:  "/tmp/csvql_history",
		VariableFile: variableFile,
		Tokens:       []Token{},
		Ast:          AST{},
		Result:       [][]string{},
		Variables:    map[string]string{},
	}

	_, err = utils.NewFile(variableFile)
	if err == nil {
		query.ReadVariable()
	}

	return query
}

func (ql *CSVQL) Execute() {
	ql.Tokenizer()
	ql.BuildAST()
	ql.ExecuteAST()
}
func (ql *CSVQL) ReadVariable() {
	file, err := os.Open(ql.VariableFile)
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
		ql.Variables[row[0]] = row[1]
	}
}
