package engine

import "fmt"

type ValueType int

const (
	TypeUndefined ValueType = iota
	TypeNull
	TypeNumber
	TypeString
	TypeBoolean
	TypeFunction
)

type Value struct {
	Type  ValueType
	Data  interface{}
	value bool
}

var Undefined = Value{Type: TypeUndefined}

func (v Value) ToString() string {
	switch v.Type {
	case TypeUndefined:
		return "undefined"
	case TypeNull:
		return "null"
	case TypeNumber:
		return fmt.Sprintf("%v", v.Data)
	case TypeString:
		return v.Data.(string)
	case TypeBoolean:
		return fmt.Sprintf("%v", v.value)
	default:
		return "[object Object]"
	}
}

func (v Value) ToNumber() float64 {
	switch v.Type {
	case TypeNumber:
		return v.Data.(float64)
	default:
		return 0
	}
}

func (v Value) ToBoolean() bool {
	switch v.Type {
	case TypeBoolean:
		return v.value
	case TypeNumber:
		return v.Data.(float64) != 0
	case TypeString:
		return v.Data.(string) != ""
	case TypeNull, TypeUndefined:
		return false
	default:
		return true
	}
}

func (v Value) Equals(other Value) bool {
	if v.Type != other.Type {
		return false
	}
	switch v.Type {
	case TypeNumber:
		return v.Data.(float64) == other.Data.(float64)
	case TypeString:
		return v.Data.(string) == other.Data.(string)
	case TypeBoolean:
		return v.value == other.value
	case TypeNull, TypeUndefined:
		return true
	default:
		return false
	}
}

func (v Value) IsFunction() bool {
	return v.Type == TypeFunction
}

func (v Value) Call(args ...Value) Value {
	if v.Type != TypeFunction {
		return Undefined
	}
	return Undefined
}
