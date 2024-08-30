package evaluator

import (
	"fmt"
	"strings"

	"github.com/ollybritton/calclang/ast"
	"github.com/ollybritton/calclang/builtins"
	"github.com/ollybritton/calclang/object"
)

// Eval evaluates a node, and returns its representation as an object.Object.
func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node, env)

	// Statements
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)

	case *ast.VariableAssignment:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}

		if isBuiltin(node.Name.Value) {
			return newError("cannot assign to builtin: %s", node.Name.Value)
		}

		if node.Name.Constant {
			err := env.SetConstant(node.Name.Value, val)
			if isError(err) {
				return err
			}
		} else {
			err := env.Set(node.Name.Value, val)
			if isError(err) {
				return err
			}
		}

		return val

	case *ast.SubroutineCall:
		expression := Eval(node.Subroutine, env)
		if isError(expression) {
			return expression
		}

		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}

		return applySubroutine(expression, args)

	// Literals
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.FloatLiteral:
		return &object.Float{Value: node.Value}

	// Expressions
	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}

		return evalPrefixExpression(node.Operator, right)

	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}

		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}

		return evalInfixExpression(left, node.Operator, right)

	case *ast.Identifier:
		return evalIdentifier(node, env)
	}

	return nil
}

func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range program.Statements {
		result = Eval(statement, env)
		fmt.Println(result.Inspect()) // TODO: probably not the right place

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}

	return result
}

func evalExpressions(exps []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object

	for _, e := range exps {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}

		result = append(result, evaluated)
	}

	return result
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	switch val := right.(type) {
	case *object.Integer:
		return &object.Integer{Value: -val.Value}
	case *object.Float:
		return &object.Float{Value: -val.Value}
	default:
		return newError("unknown operator: -%s", right.Type())
	}
}

func evalInfixExpression(left object.Object, operator string, right object.Object) object.Object {
	left, right = coerceInfix(left, operator, right)

	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(left, operator, right)

	case left.Type() == object.FLOAT_OBJ && right.Type() == object.FLOAT_OBJ:
		return evalFloatInfixExpression(left, operator, right)

	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())

	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIntegerInfixExpression(left object.Object, operator string, right object.Object) object.Object {
	leftInt := left.(*object.Integer)
	rightInt := right.(*object.Integer)

	switch operator {
	case "+":
		return &object.Integer{Value: leftInt.Value + rightInt.Value}
	case "-":
		return &object.Integer{Value: leftInt.Value - rightInt.Value}
	case "*":
		return &object.Integer{Value: leftInt.Value * rightInt.Value}
	case "/":
		if rightInt.Value == 0 {
			return newError("division error: division by zero")
		}

		if leftInt.Value%rightInt.Value == 0 {
			return &object.Integer{Value: leftInt.Value / rightInt.Value}
		}

		lf := object.IntegerToFloat(leftInt)
		rf := object.IntegerToFloat(rightInt)

		return &object.Float{Value: lf.Value / rf.Value}
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalFloatInfixExpression(left object.Object, operator string, right object.Object) object.Object {
	lf := left.(*object.Float)
	rf := right.(*object.Float)

	switch operator {
	case "+":
		return &object.Float{Value: lf.Value + rf.Value}
	case "-":
		return &object.Float{Value: lf.Value - rf.Value}
	case "*":
		return &object.Float{Value: lf.Value * rf.Value}
	case "/":
		if rf.Value == 0 {
			return newError("division error: division by zero")
		}

		return &object.Float{Value: lf.Value / rf.Value}
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}

	if builtin, ok := builtins.Builtins[strings.ToUpper(node.Value)]; ok {
		return builtin
	}

	return newError("identifier not found: " + node.Value)
}

func applySubroutine(sub object.Object, args []object.Object) object.Object {
	switch sub := sub.(type) {

	case *object.Builtin:
		return sub.Fn(args...)

	default:
		return newError("not a subroutine, function or builtin: %s", sub.Type())
	}
}
