package main

import (
	"os"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/alecthomas/repr"
)

var (
	calcLexer = lexer.MustSimple([]lexer.SimpleRule{
		{Name: `Ident`, Pattern: `[ABCDEFxyM]`},
		{Name: `Float`, Pattern: `\d+(?:\.\d+)?`},
		{Name: `Arrow`, Pattern: `->`},
		{Name: `Colon`, Pattern: `:`},
		{Name: `whitespace`, Pattern: `\s+`},
	})

	parser = participle.MustBuild[CALC](
		participle.Lexer(calcLexer),
	)
)

type CALC struct {
	Statements []*Statement `@@*`
}

type Statement struct {
	Num   float64 `@Float Arrow`
	Ident string  `@Ident Colon?`
}

func main() {
	ini, err := parser.Parse("", os.Stdin)
	if err != nil {
		panic(err)
	}
	repr.Println(ini, repr.Indent("\t"))
}
