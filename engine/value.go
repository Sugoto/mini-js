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
	TypeObject
	TypeReturn
)

type Value struct {
	Type       ValueType
	Data       interface{}
	value      bool
	Properties map[string]Value
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
		if v.Data != nil {
			return fmt.Sprintf("%v", v.Data.(bool))
		}
		return fmt.Sprintf("%v", v.value)
	case TypeFunction:
		return "[Function]"
	case TypeObject:
		return "[object Object]"
	case TypeReturn:
		if ret, ok := v.Data.(*ReturnValue); ok {
			return ret.Value.ToString()
		}
		return "undefined"
	default:
		return "[object Object]"
	}
}

func (v Value) ToNumber() float64 {
	switch v.Type {
	case TypeNumber:
		return v.Data.(float64)
	case TypeString:
		return 0
	case TypeBoolean:
		if v.Data != nil {
			if v.Data.(bool) {
				return 1
			}
		} else if v.value {
			return 1
		}
		return 0
	default:
		return 0
	}
}

func (v Value) ToBoolean() bool {
	switch v.Type {
	case TypeBoolean:
		if v.Data != nil {
			return v.Data.(bool)
		}
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
		if v.Data != nil && other.Data != nil {
			return v.Data.(bool) == other.Data.(bool)
		}
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

func (v Value) GetProperty(name string) Value {
	if v.Type == TypeObject && v.Properties != nil {
		if prop, ok := v.Properties[name]; ok {
			return prop
		}
	}

	if v.Type == TypeObject && v.Data != nil {
		if v.Data.(string) == "console" && name == "log" {
			return Value{
				Type: TypeFunction,
				Data: &ConsoleLogFunction{},
			}
		}
	}

	return Undefined
}

func (v Value) SetProperty(name string, value Value) {
	if v.Type != TypeObject {
		return
	}

	if v.Properties == nil {
		v.Properties = make(map[string]Value)
	}

	v.Properties[name] = value
}

type Function struct {
	Parameters []*Identifier
	Body       *BlockStatement
	Env        *Environment
}

type ConsoleLogFunction struct{}

func (clf *ConsoleLogFunction) Call(args ...Value) Value {
	for _, arg := range args {
		fmt.Println(arg.ToString())
	}
	return Undefined
}

func (v Value) Call(args ...Value) Value {
	if v.Type != TypeFunction {
		return Undefined
	}

	// Handle built-in functions like console.log
	if consoleLog, ok := v.Data.(*ConsoleLogFunction); ok {
		return consoleLog.Call(args...)
	}

	fn, ok := v.Data.(*Function)
	if !ok {
		return Undefined
	}

	newEnv := NewEnvironment()
	newEnv.store = fn.Env.store

	for i, param := range fn.Parameters {
		if i < len(args) {
			newEnv.Set(param.Value, args[i])
		} else {
			newEnv.Set(param.Value, Undefined)
		}
	}

	tempInterpreter := &Interpreter{env: newEnv}

	var result Value = Undefined
	for _, stmt := range fn.Body.Statements {
		result = tempInterpreter.evalStatement(stmt)
		if ret, ok := stmt.(*ReturnStatement); ok {
			return tempInterpreter.evalExpression(ret.ReturnValue)
		}
	}

	return result
}

func (v Value) Add(other Value) Value {
	if v.Type == TypeNumber && other.Type == TypeNumber {
		result := v.Data.(float64) + other.Data.(float64)
		return Value{
			Type: TypeNumber,
			Data: result,
		}
	}
	if v.Type == TypeString || other.Type == TypeString {
		return Value{
			Type: TypeString,
			Data: v.ToString() + other.ToString(),
		}
	}
	return Undefined
}

func (v Value) Subtract(other Value) Value {
	if v.Type == TypeNumber && other.Type == TypeNumber {
		result := v.Data.(float64) - other.Data.(float64)
		return Value{
			Type: TypeNumber,
			Data: result,
		}
	}
	return Undefined
}

func (v Value) Multiply(other Value) Value {
	if v.Type == TypeNumber && other.Type == TypeNumber {
		result := v.Data.(float64) * other.Data.(float64)
		return Value{
			Type: TypeNumber,
			Data: result,
		}
	}
	return Undefined
}

func (v Value) Divide(other Value) Value {
	if v.Type == TypeNumber && other.Type == TypeNumber {
		if other.Data.(float64) == 0 {
			return Undefined
		}
		result := v.Data.(float64) / other.Data.(float64)
		return Value{
			Type: TypeNumber,
			Data: result,
		}
	}
	return Undefined
}
