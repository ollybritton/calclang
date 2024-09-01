package repl

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/alecthomas/repr"
	"github.com/c-bata/go-prompt"
	au "github.com/logrusorgru/aurora"
	"github.com/ollybritton/calclang/evaluator"
	"github.com/ollybritton/calclang/lexer"
	"github.com/ollybritton/calclang/object"
	"github.com/ollybritton/calclang/parser"
	"github.com/ollybritton/calclang/token"
)

// Repl represents a repl. It can be used to lex, parse or evaluate calclang code.
type Repl struct {
	Buffer bytes.Buffer
	Prompt *prompt.Prompt

	Mode  string // Either "lex", "parse" or "eval"
	Level int

	Env *object.Environment
}

// New returns a new, initialised REPL.
func New() *Repl {
	r := &Repl{Mode: "eval"}
	r.Prompt = prompt.New(
		r.Execute,
		r.Completor,

		prompt.OptionLivePrefix(r.Prefix),
		prompt.OptionTitle("calclang"),

		prompt.OptionSuggestionBGColor(prompt.Red),
		prompt.OptionSuggestionTextColor(prompt.Black),

		prompt.OptionSelectedSuggestionBGColor(prompt.Red),
		prompt.OptionSelectedDescriptionBGColor(prompt.Turquoise),
		prompt.OptionSelectedDescriptionTextColor(prompt.Black),
		prompt.OptionSelectedSuggestionTextColor(prompt.Black),

		prompt.OptionInputTextColor(prompt.Turquoise),
	)

	return r
}

// Execute is what executes a command inside the REPL.
func (r *Repl) Execute(input string) {
	r.Level = 0

	if strings.HasPrefix(input, "%") {
		switch input[1:] {
		case "lex", "tokenize", "split":
			r.Mode = "lex"
			fmt.Println(au.Green("Mode set to 'lex'."))
			fmt.Println("")

			return

		case "parse", "ast":
			r.Mode = "parse"
			fmt.Println(au.Green("Mode set to 'parse'."))
			fmt.Println("")

			return

		case "eval", "exec":
			r.Mode = "eval"
			fmt.Println(au.Green("Mode set to 'eval'."))
			fmt.Println("")

			return

		case "buf":
			input = Buffer(false)

			if input != "" {
				fmt.Println("")
				fmt.Println(au.Green("Added from buffer:").Bold())
				fmt.Println(au.Yellow(input))
			}

		case "clearbuf":
			input = Buffer(true)

			if input != "" {
				fmt.Println("")
				fmt.Println(au.Green("Added from buffer:").Bold())
				fmt.Println(au.Yellow(input))
			}

		case "help":
			Help()
			return

		default:
			message := au.Red(au.Bold(
				fmt.Sprintf("No magic command %q found.", input),
			))

			fmt.Println(message)
			fmt.Println("")

			return
		}
	}

	switch r.Mode {
	case "lex":
		r.Lex(input)
	case "parse":
		r.Parse(input)
	case "eval":
		r.Eval(input)
	}
}

// Completor is what completes input inside the REPL.
func (r *Repl) Completor(input prompt.Document) []prompt.Suggest {
	w := input.GetWordBeforeCursor()
	if w == "" {
		return []prompt.Suggest{}
	}

	return prompt.FilterHasPrefix(suggestions, w, true)
}

// Prefix is what calculates the prefix/identation level.
func (r *Repl) Prefix() (string, bool) {
	if r.Level > 0 {
		indent := strings.Repeat("  ", r.Level)
		return "(" + r.Mode + ") > " + indent, true
	}

	return "(" + r.Mode + ") > ", true
}

// Lex lexes a given input string, and displays the results to stdout.
func (r *Repl) Lex(input string) {
	l := lexer.New(input)

	tokens := []token.Token{}
	tok := l.NextToken()

	for tok.Type != token.EOF {
		tokens = append(tokens, tok)
		tok = l.NextToken()
	}

	fmt.Println("")

	for i, t := range tokens {
		num := au.Blue(fmt.Sprintf("[%d]", i))
		fmt.Printf("%v %v\n", num, PrettyToken(t))
	}

	fmt.Println("")
}

// Parse parses a given input string, and displays the results to stdout.
func (r *Repl) Parse(input string) {
	l := lexer.New(input)
	p := parser.New(l)

	program := p.Parse()
	if len(p.Errors()) != 0 {
		Errors(p.Errors())
	}

	fmt.Println(program.String())
	fmt.Println(au.Faint(repr.String(program, repr.Indent("\t"))))
	fmt.Println("")
}

// Eval evaluates a given input string, and displays the results to stdout.
func (r *Repl) Eval(input string) {
	obj, errors := evaluator.EvalString(input, r.Env)

	if len(errors) != 0 {
		Errors(errors)
	}

	if obj == nil {
		return
	}

	if obj.Type() == object.ERROR_OBJ {
		fmt.Println(au.Red(obj.Inspect()).Bold())
		return
	}

	fmt.Println(au.Green(obj.Inspect()))
	fmt.Println("")
}

// Start starts the REPL.
func (r *Repl) Start() {
	Info()

	for {
		input := r.Prompt.Input()

		switch {
		case input == "exit" || input == "quit", input == "%exit" || input == "%quit":
			os.Exit(0)
		case input == "ping":
			fmt.Println("pong")
			fmt.Println("")
		default:
			r.Execute(input)
		}
	}
}
