package builtins

import (
	"math"
	"math/rand"

	"github.com/ollybritton/calclang/object"
)

const FLOAT_EQUALITY_TOL float64 = 1e-6

// BuiltinRandomInt generates a random integer object between two bounds.
func BuiltinRandomInt(args ...object.Object) object.Object {
	if len(args) != 2 {
		return newError("wrong number of arguments. got=%d, want=2", len(args))
	}

	lower, ok := args[0].(*object.Integer)
	if !ok {
		return newError("argument 1 to `RANDOM_INT` not supported, got=%s", args[0].Type())
	}

	upper, ok := args[1].(*object.Integer)
	if !ok {
		return newError("argument 2 to `RANDOM_INT` not supported, got=%s", args[1].Type())
	}

	val := rand.Intn(int(upper.Value-lower.Value+1)) + int(lower.Value)
	return &object.Integer{Value: int64(val)}
}

// BuiltinFloor will floor a float. It has no effect on integers.
func BuiltinFloor(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1", len(args))
	}

	switch val := args[0].(type) {
	case *object.Float:
		return &object.Integer{Value: int64(math.Floor(val.Value))}
	case *object.Integer:
		return val
	default:
		return newError("argument to `FLOOR` not supported, got=%s", args[0].Type())
	}
}

// BuiltinRound will round a float. It has no effect on integers.
func BuiltinRound(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1", len(args))
	}

	switch val := args[0].(type) {
	case *object.Float:
		return &object.Integer{Value: int64(math.Round(val.Value))}
	case *object.Integer:
		return val
	default:
		return newError("argument to `ROUND` not supported, got=%s", args[0].Type())
	}
}

// BuiltinCeil will round a float up. It has no effect on integers.
func BuiltinCeil(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1", len(args))
	}

	switch val := args[0].(type) {
	case *object.Float:
		return &object.Integer{Value: int64(math.Ceil(val.Value))}
	case *object.Integer:
		return val
	default:
		return newError("argument to `CEIL` not supported, got=%s", args[0].Type())
	}
}

// BuiltinSqrt will find the square root of an integer or a float.
func BuiltinSqrt(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1", len(args))
	}

	var inp float64

	switch val := args[0].(type) {
	case *object.Float:
		inp = val.Value
	case *object.Integer:
		inp = float64(val.Value)
	default:
		return newError("argument to `SQRT` not supported, got=%s", args[0].Type())
	}

	result := math.Sqrt(float64(inp))

	if math.IsNaN(result) {
		return newError("MathERROR")
	}

	return &object.Float{Value: result}
}

// BuiltinDelta
func BuiltinKronDelta(args ...object.Object) object.Object {
	if len(args) == 0 {
		return newError("wrong number of arguments. got=0, want:>=1")
	}

	floatArgs := []float64{}

	for _, arg := range args {
		switch arg := arg.(type) {
		case *object.Float:
			floatArgs = append(floatArgs, arg.Value)
		case *object.Integer:
			floatArgs = append(floatArgs, float64(arg.Value))
		}
	}

	start := floatArgs[0]
	same := true

	for _, arg := range floatArgs {
		if math.Abs(arg-start) > FLOAT_EQUALITY_TOL {
			same = false
			break
		}
	}

	if same {
		return &object.Integer{Value: 1}
	}

	return &object.Integer{Value: 0}
}
