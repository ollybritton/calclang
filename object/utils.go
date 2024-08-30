package object

import (
	"math"
)

// IntegerToFloat converts an integer object to a float object.
func IntegerToFloat(i *Integer) *Float {
	return &Float{Value: float64(i.Value)}
}

// FloatToInteger converts a flaot object to an integer object.
func FloatToInteger(f *Float) *Integer {
	val := int64(math.Round(f.Value))
	return &Integer{Value: val}
}
