package builtins

import (
	"fmt"

	"github.com/ollybritton/calclang/object"
)

// Builtins maps the name of a builtin function within the program to the actual function.
var Builtins = make(map[string]*object.Builtin)

func newError(message string, args ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(message, args...)}
}

func init() {
	Builtins["RANDOM_INT"] = &object.Builtin{Fn: BuiltinRandomInt, Strict: true}
	Builtins["ROUND"] = &object.Builtin{Fn: BuiltinRound, Strict: true}
	Builtins["SQRT"] = &object.Builtin{Fn: BuiltinSqrt, Strict: true}

	Builtins["P"] = &object.Builtin{Fn: BuiltinPrint, Strict: false}
	Builtins["DELTA"] = &object.Builtin{Fn: BuiltinKronDelta, Strict: false}
	Builtins["FLOOR"] = &object.Builtin{Fn: BuiltinFloor, Strict: false}
	Builtins["CEIL"] = &object.Builtin{Fn: BuiltinCeil, Strict: false}
}
