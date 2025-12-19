package pkg

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jessevdk/go-flags"
)

type Options struct {
	// Slice of bool will append 'true' each time the option
	// is encountered (can be set multiple times, like -vvv)
	Verbose []bool `short:"v" long:"verbose" description:"Show verbose debug information"`

	// Example of automatic marshalling to desired type (uint)
	Offset uint `long:"offset" description:"Offset"`

	// Example of a callback, called each time the option is found.
	Call func(string) `short:"c" description:"Call phone number"`

	// Example of a required flag
	Name string `short:"n" long:"name" description:"A name" required:"true"`

	// Example of a flag restricted to a pre-defined set of strings
	Animal string `long:"animal" choice:"cat" choice:"dog"`

	// Example of a value name
	File string `short:"f" long:"file" description:"A file" value-name:"FILE"`

	// Example of a pointer
	Ptr *int `short:"p" description:"A pointer to an integer"`

	// Example of a slice of strings
	StringSlice []string `short:"s" description:"A slice of strings"`

	// Example of a slice of pointers
	PtrSlice []*string `long:"ptrslice" description:"A slice of pointers to string"`

	// Example of a map
	IntMap map[string]int `long:"intmap" description:"A map from string to int"`

	// Example of env variable
	Thresholds []int `long:"thresholds" default:"1" default:"2" env:"THRESHOLD_VALUES"  env-delim:","`
}

type VariableOptions struct {
	VariableMap map[string]string `short:"s" long:"set" description:"Set variable" key-value-delimiter:"=" `
	Get         string            `short:"g" long:"get" description:"Get variable"`
	Unset       string            `short:"u" long:"unset" description:"Unset variable"`
}

func FlagParser(input string) {

	var options VariableOptions

	// args := []string{
	// 	"-vv",
	// 	"--offset=5",
	// 	"-n", "Me",
	// 	"--animal", "dog", // anything other than "cat" or "dog" will raise an error
	// 	"-p", "3",
	// 	"-s", "hello",
	// 	"-s", "world",
	// 	"--ptrslice", "hello",
	// 	"--ptrslice", "world",
	// 	"--intmap", "a:1",
	// 	"--intmap", "b:5",
	// 	"arg1",
	// 	"arg2",
	// 	"arg3",
	// }
	args := []string{
		// "-s",
		// "name=luannguyen",

		"-g", "name",
	}
	var parser = flags.NewParser(&options, flags.Default)
	parseArgs, err := parser.ParseArgs(args)
	PrintPretty("options", options)
	PrintPretty("parseArgs", parseArgs)
	PrintPretty("error", err)

	// set name=luannguyen
	// set -l
	// set -g name
	// set -u name
	// set -r

}

func (ql *CSVQL) ReplVariable(line string) {
	FlagParser(line)
	expression, _ := strings.CutPrefix(line, "variable")

	switch {
	case strings.Contains(expression, "-l"):
		{
			ql.GetVariables()
		}
	case strings.Contains(expression, "-r"):
		{
			ql.GetVariables()
		}
	case strings.Contains(expression, "-g"):
		{
			expression = strings.ReplaceAll(expression, " ", "")
			variable, _ := strings.CutPrefix(expression, "-g")
			value, ok := ql.Variables[variable]
			if !ok {
				log.Println("Invalid variable. Please try again.")
				return
			}
			log.Println(value)
		}

	case strings.Contains(expression, "-u"):
		{
			ql.GetVariables()
		}
	default:
		{
			expression = strings.ReplaceAll(expression, " ", "")

			expressionArr := strings.Split(expression, "=")
			if len(expressionArr) != 2 {
				fmt.Println("Set variable failed. Please try again.")
				return
			}

			variable := expressionArr[0]
			value := expressionArr[1]
			// ql.Set(expressionArr[0], expressionArr[1])
			ql.Variables[variable] = value
			row := fmt.Sprintf("%v,%v\n", variable, value)
			file, fileErr := os.OpenFile(ql.VariableFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			fileInfo, _ := file.Stat()
			PrintPretty("fileInfo", fileInfo)
			PrintPretty("fileErr", fileErr)
			defer file.Close()
			_, err := file.WriteString(row)
			PrintPretty("err", err)
		}
	}

}
