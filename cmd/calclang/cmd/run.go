package cmd

import (
	"fmt"
	"os"

	"github.com/ollybritton/calclang/evaluator"
	"github.com/ollybritton/calclang/object"
	"github.com/ollybritton/calclang/repl"

	au "github.com/logrusorgru/aurora"
	"github.com/ollybritton/calclang/lexer"
	"github.com/ollybritton/calclang/parser"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run [filename]",
	Args:  cobra.MaximumNArgs(1),
	Short: "run runs a .calc file and displays the output",
	Long: `run will run a file containing calclang code.
For now, it will also print the result of the evaluation.`,
	Run: func(cmd *cobra.Command, args []string) {
		command, err := cmd.Flags().GetString("command")
		if err != nil {
			fmt.Println(au.Bold(au.Red("Could not fetch flag:")))
			fmt.Println(au.Red(err))
		}

		var str string

		if command != "" {
			str = command
		} else {
			bytes, err := os.ReadFile(args[0])
			if err != nil {
				fmt.Println(au.Bold(au.Red("Could not read file:")))
				fmt.Println(au.Red(err))
			}

			str = string(bytes)
		}

		l := lexer.New(str)
		p := parser.New(l)

		program := p.Parse()
		if len(p.Errors()) != 0 {
			repl.Errors(p.Errors())
			return
		}

		eval := evaluator.Eval(program, object.NewEnvironment())
		if eval == nil {
			return
		}

		if eval.Type() == object.ERROR_OBJ {
			fmt.Println(au.Red(eval.Inspect()).Bold())
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	runCmd.Flags().StringP("command", "c", "", "Command to run before exiting")
}
