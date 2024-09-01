package builtins

import (
	"fmt"

	"github.com/ollybritton/calclang/object"
)

// BuiltinPrint will print the value of an expression and return that same expression.
func BuiltinPrint(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1", len(args))
	}

	fmt.Println(args[0].Inspect())
	return args[0]
}
