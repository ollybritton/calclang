package evaluator

import "github.com/ollybritton/calclang/object"

// coerceInfix will convert the types in an infix expression automatically so that objects can be used with one another
// without having to deal with converting types yourself.
//
// Rules:
// int, float => float & float
// float, int => float & float
func coerceInfix(left object.Object, operator string, right object.Object) (object.Object, object.Object) {
	_ = operator

	switch {
	case left.Type() == object.FLOAT_OBJ && right.Type() == object.INTEGER_OBJ:
		x := left.(*object.Float)
		y := right.(*object.Integer)

		return x, object.IntegerToFloat(y)

	case left.Type() == object.INTEGER_OBJ && right.Type() == object.FLOAT_OBJ:
		x := left.(*object.Integer)
		y := right.(*object.Float)

		return object.IntegerToFloat(x), y

	default:
		return left, right
	}
}
