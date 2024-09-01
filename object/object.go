package object

import (
	"fmt"
	"strconv"
)

// Type represents a type of object, such as an integer or a subroutine.
type Type string

// Definition of object types.
const (
	INTEGER_OBJ      = "INTEGER"
	FLOAT_OBJ        = "FLOAT"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
	BUILTIN_OBJ      = "BUILTIN"
	ERROR_OBJ        = "ERROR"
)

// Object is an interface which allows different objects to be represented.
type Object interface {
	Type() Type      // Type reveals an object's type
	Inspect() string // Inspect gets the value of the object as a string.
}

// BuiltinFunction represents an external function that is avaliable inside an calclang
// program.
type BuiltinFunction func(args ...Object) Object

// Builtin represents a builtin inside the program.
type Builtin struct {
	Fn     BuiltinFunction
	Strict bool
}

func (b *Builtin) Type() Type      { return BUILTIN_OBJ }
func (b *Builtin) Inspect() string { return "<builtin>" }

// Integer represents an integer within the program.
type Integer struct {
	Value int64
}

func (i *Integer) Type() Type      { return INTEGER_OBJ }
func (i *Integer) Inspect() string { return fmt.Sprintf("%d", i.Value) }

// Float represents an Float within the program.
type Float struct {
	Value float64
}

func (f *Float) Type() Type      { return FLOAT_OBJ }
func (f *Float) Inspect() string { return strconv.FormatFloat(f.Value, 'f', -1, 64) }

// ReturnValue represents a value that is being returned from a subroutine or from a program as a whole.
type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() Type      { return RETURN_VALUE_OBJ }
func (rv *ReturnValue) Inspect() string { return rv.Value.Inspect() }

// Error represents an error that occurs during the evalutation of the programming language.
type Error struct {
	Message string
}

func (e *Error) Type() Type      { return ERROR_OBJ }
func (e *Error) Inspect() string { return "ERROR: " + e.Message }
