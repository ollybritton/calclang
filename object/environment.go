package object

import (
	"fmt"
	"math"
)

// Environment represents the variables and identifiers inside their program, mapped to their actual Object values.
type Environment struct {
	store     map[string]Object
	constants map[string]Object
	outer     *Environment
}

// NewEnvironment creates a new environment.
func NewEnvironment() *Environment {
	s := make(map[string]Object)
	c := make(map[string]Object)

	c["pi"] = &Float{Value: math.Pi}
	c["e"] = &Float{Value: math.E}

	return &Environment{store: s, constants: c, outer: nil}
}

// NewEnclosedEnvironment creates a new enclosed environment, extending from a previous.
func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer

	return env
}

// Get gets an object by name.
func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if ok {
		// In normal store
		return obj, ok
	}

	obj, ok = e.constants[name]
	if ok {
		// In constant store
		return obj, ok
	}

	if e.outer != nil {
		// Inside outer environment
		obj, ok = e.outer.Get(name)
		if ok {
			return obj, ok
		}
	}

	// return obj, ok
	return nil, false
}

// Set sets an object by name.
func (e *Environment) Set(name string, value Object) Object {

	if _, ok := e.constants[name]; ok {
		return &Error{Message: fmt.Sprintf("cannot assign to constant %s", name)}
	}

	e.store[name] = value
	return value
}

// Keys gets the list of all symbols.
func (e *Environment) Keys() map[string]bool {
	symbols := make(map[string]bool)

	for k := range e.store {
		symbols[k] = true
	}

	for k := range e.constants {
		symbols[k] = true
	}

	return symbols
}
