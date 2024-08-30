package evaluator

import (
	"errors"
	"io"
	"os"

	"github.com/ollybritton/calclang/lexer"
	"github.com/ollybritton/calclang/object"
	"github.com/ollybritton/calclang/parser"
)

// EvalString will execute a string of calclang code.
func EvalString(str string, env *object.Environment) (object.Object, []error) {
	l := lexer.New(str)
	p := parser.New(l)

	program := p.Parse()
	if len(p.Errors()) != 0 {
		return nil, p.Errors()
	}

	eval := Eval(program, env)
	if eval == nil {
		return nil, []error{}
	}

	if eval.Type() == object.ERROR_OBJ {
		return nil, []error{errors.New(eval.Inspect())}
	}

	return eval, []error{}
}

// EvalFile will execute a file containing calclang code.
func EvalFile(f *os.File, env *object.Environment) (object.Object, []error) {
	bytes, err := io.ReadAll(f)
	if err != nil {
		return nil, []error{err}
	}

	return EvalString(string(bytes), env)
}
