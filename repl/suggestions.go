package repl

import "github.com/c-bata/go-prompt"

var suggestions = []prompt.Suggest{
	{Text: "%help", Description: "Print help text."},

	{Text: "%lex", Description: "Put the REPL into lex mode."},
	{Text: "%parse", Description: "Put the REPL into parse mode."},
	{Text: "%eval", Description: "Put the REPL into eval mode."},

	{Text: "exit", Description: "Exit the REPL."},
	{Text: "quit", Description: "Exit the REPL."},
}
