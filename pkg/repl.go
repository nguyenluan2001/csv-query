package pkg

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	"github.com/chzyer/readline"
	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/tw"
)

func (ql *CSVQL) Render() {
	table := tablewriter.NewWriter(os.Stdout)

	for i, row := range ql.Result {
		if i == 0 {
			header := append([]string{"#"}, row...)
			table.Header(header)
		} else {
			row = append([]string{fmt.Sprintf("%d", i)}, row...)
			table.Append(row)
		}
	}

	table.Render()
	fmt.Println(fmt.Sprintf("(%d rows returned. Executed in %vms)", len(ql.Result)-1, ql.Duration))
}

func NewTable(header []string, rows [][]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.Header(header)
	table.Bulk(rows)
	table.Render()
}

// === Command handler ===
func (ql *CSVQL) Export(filePath string) {
	if !strings.HasPrefix(filePath, "/") {
		filePath = path.Join(ql.DatabasePath, filePath)
	}
	filePath = strings.ReplaceAll(filePath, " ", "")
	fmt.Println("filePath", filePath)

	_, err := os.Stat(filePath)

	if err == nil {
		fmt.Println("File existed. Please try again.")
		return
	}

	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Something went wrong. Please try again.")
		return
	}

	defer file.Close()
	for _, row := range ql.Result {
		rowStr := strings.Join(row, ",")
		file.WriteString(fmt.Sprintf("%v\n", rowStr))
	}

}

func (ql *CSVQL) Tables() {
	entries, err := os.ReadDir(ql.DatabasePath)
	table := tablewriter.NewWriter(os.Stdout)

	if err != nil {
		fmt.Println("Get list table failed. Please try again.")
		return
	}

	rows := [][]string{}
	for i, entry := range entries {
		row := append([]string{fmt.Sprintf("%v", i+1)}, entry.Name())
		rows = append(rows, row)
	}
	table.Header([]string{"#", "Table"})
	table.Bulk(rows)
	table.Render()
	fmt.Println(fmt.Sprintf("(%d tables)", len(entries)))

}

func (ql *CSVQL) FetchHistory() [][]string {
	file, err := os.Open(ql.HistoryFile)
	result := [][]string{}

	if err != nil {
		fmt.Println("Get history failed. Please try again.")
		return result
	}

	reader := bufio.NewReader(file)
	countRow := 1
	for {
		row, _, rowErr := reader.ReadLine()
		if rowErr == io.EOF {
			break // End of file reached
		}
		if rowErr != nil {
			// Handle error during reading
		}
		newRow := append([]string{fmt.Sprintf("%v", countRow)}, []string{string(row)}...)
		result = append(result, newRow)
		countRow++
	}
	return result
}

func (ql *CSVQL) History(line string) {
	if strings.Contains(line, "-d") {
		fmt.Println("Clear histoy")
		cmd := exec.Command("sh", "-c", fmt.Sprintf("cat /dev/null > %v", ql.HistoryFile))
		cmd.Output()
		ql.Readline.ResetHistory()
		ql.Readline.Operation.ResetHistory()
		return
	}
	result := ql.FetchHistory()

	table := tablewriter.NewTable(os.Stdout,
		tablewriter.WithConfig(tablewriter.Config{
			Row: tw.CellConfig{
				Formatting:   tw.CellFormatting{AutoWrap: tw.WrapNormal}, // Wrap long content
				Alignment:    tw.CellAlignment{Global: tw.AlignLeft},     // Left-align rows
				ColMaxWidths: tw.CellWidth{Global: 100},
			},
			Footer: tw.CellConfig{
				Alignment: tw.CellAlignment{Global: tw.AlignRight},
			},
		}),
	)
	table.Header([]string{"#", "Command"})
	table.Bulk(result)
	table.Render()
	// PrintPretty("result", result)

}

func (ql *CSVQL) ReuseHistory(number int) string {
	history := ql.FetchHistory()
	cmd := ""

	for _, item := range history {
		if Stringify(item[0]) == Stringify(number) {
			cmd = item[1]
			break
		}
	}
	return cmd
}

func (ql *CSVQL) Help() {
	cmds := [][]string{
		[]string{"exit", "Exit"},
		[]string{"clear", "Clear screen"},
		[]string{"export", "Export result to file"},
		[]string{"tables", "Get list table"},
		[]string{"history", "Display or manipulate the history list"},
	}
	fmt.Println("Commands:")
	for _, cmd := range cmds {
		cmdName := cmd[0]
		if strings.Contains(cmd[0], "help") {
			continue
		}
		desc := cmd[1]
		fmt.Println(fmt.Sprintf("   %v: %v", cmdName, desc))
	}
}

func (ql *CSVQL) Set(variable, value string) {
	ql.Variables[variable] = value
}

func (ql *CSVQL) GetVariables() {
	rows := [][]string{}
	count := 1
	for key, value := range ql.Variables {
		rows = append(rows, []string{Stringify(count), key, value})
		count++
	}
	NewTable([]string{"#", "Name", "Value"}, rows)
}

func SetPrompt(rl *readline.Instance, cmd string) {
	pattern := `csvql> %v`
	prompt := fmt.Sprintf(pattern, cmd)
	rl.SetPrompt(prompt)

}

var completer = readline.NewPrefixCompleter(
	readline.PcItem("exit",
		readline.PcItem("Exit."),
	),
	readline.PcItem("clear",
		readline.PcItem("Clear screen."),
	),
	readline.PcItem("export",
		readline.PcItem("Export result to file."),
	),
	readline.PcItem("tables",
		readline.PcItem("Get list table."),
	),
	readline.PcItem("history",
		readline.PcItem("Display or manipulate the history list."),
	),
	readline.PcItem("set",
		readline.PcItem("Set global variable."),
	),
	readline.PcItem("help"),
)

func (ql *CSVQL) Repl() {
	rl, err := readline.NewEx(&readline.Config{
		Prompt:                 "csvql> ",
		HistoryFile:            ql.HistoryFile,
		DisableAutoSaveHistory: true,
		AutoComplete:           completer,
	})
	if err != nil {
		panic(err)
	}
	defer rl.Close()
	ql.Readline = rl

	var cmds []string
	for {
		line, err := rl.Readline()
		if err != nil {
			break
		}
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		switch {
		case line == "exit":
			{
				rl.SaveHistory(line)
				os.Exit(0)
			}
		case line == "help":
			{
				rl.SaveHistory(line)
				ql.Help()
			}
		case line == "tables":
			{
				rl.SaveHistory(line)
				ql.Tables()
			}
		case strings.HasPrefix(line, "history"):
			{
				rl.SaveHistory(line)
				ql.History(line)
			}
		case strings.HasPrefix(line, "!"):
			{
				numberStr, _ := strings.CutPrefix(line, "!")
				number, err := StringToInt(numberStr)

				if err != nil {
					fmt.Println(fmt.Sprintf("event not found: %v", numberStr))
					return
				}

				cmd := ql.ReuseHistory(number)

				if len(cmd) == 0 {
					return
				}

				_, err = rl.WriteStdin([]byte(cmd))

				if err != nil {
					fmt.Println("Failed to execute command. Please try again.")
					return
				}

			}
		case line == "clear":
			{
				rl.SaveHistory(line)
				cmd := exec.Command("clear")
				cmd.Stdout = os.Stdout
				cmd.Run()
				rl.Clean()
				rl.SetPrompt("csvql> ")
			}
		case strings.HasPrefix(line, "export"):
			{
				rl.SaveHistory(line)
				fileName, _ := strings.CutPrefix(line, "export")
				ql.Export(fileName)
			}
		case strings.HasSuffix(line, `\`):
			{

				cmds = append(cmds, line[0:len(line)-1])
				rl.SetPrompt("> ")
			}
		case strings.HasPrefix(line, "set"):
			{
				expression, _ := strings.CutPrefix(line, "set")
				expression = strings.ReplaceAll(expression, " ", "")

				expressionArr := strings.Split(expression, "=")
				if len(expressionArr) != 2 {
					fmt.Println("Set variable failed. Please try again.")
					return
				}
				ql.Set(expressionArr[0], expressionArr[1])

			}
		case line == "variables":
			{
				ql.GetVariables()
			}
		default:
			{
				cmds = append(cmds, line)
				cmd := strings.Join(cmds, " ")

				rl.SaveHistory(cmd)

				rl.SetPrompt("csvql> ")
				ql.Sql = cmd
				start := time.Now()
				ql.Execute()
				duration := time.Since(start)
				ql.Duration = float64(duration.Milliseconds())
				ql.Render()

				cmds = cmds[:0]
			}
		}
	}

}
