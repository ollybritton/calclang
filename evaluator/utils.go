package evaluator

import (
	"fmt"
	"strings"

	"github.com/ollybritton/calclang/builtins"
	"github.com/ollybritton/calclang/object"
)

func newError(message string, args ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(message, args...)}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}

	return false
}

func isBuiltin(name string) bool {
	if _, ok := builtins.Builtins[strings.ToUpper(name)]; ok {
		return true
	}

	return false
}
