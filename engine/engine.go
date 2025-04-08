package engine

type Environment struct {
	store map[string]Value
	outer *Environment
}

func NewEnvironment() *Environment {
	return &Environment{
		store: make(map[string]Value),
		outer: nil,
	}
}

func (e *Environment) Get(name string) (Value, bool) {
	val, ok := e.store[name]
	if !ok && e.outer != nil {
		return e.outer.Get(name)
	}
	return val, ok
}

func (e *Environment) Set(name string, val Value) {
	e.store[name] = val
}

type Interpreter struct {
	env *Environment
}

func NewInterpreter() *Interpreter {
	return &Interpreter{
		env: NewEnvironment(),
	}
}

func (i *Interpreter) SetGlobal(name string, value interface{}) error {
	var v Value
	switch val := value.(type) {
	case Value:
		v = val
	case float64:
		v = Value{Type: TypeNumber, Data: val}
	case string:
		v = Value{Type: TypeString, Data: val}
	case bool:
		v = Value{Type: TypeBoolean, value: val}
	default:
		v = Undefined
	}
	i.env.Set(name, v)
	return nil
}

func (i *Interpreter) Eval(code string) (Value, error) {
	lexer := NewLexer(code)
	parser := NewParser(lexer)
	program := parser.ParseProgram()

	return i.evalProgram(program), nil
}

func (i *Interpreter) evalProgram(program *Program) Value {
	var result Value

	for _, statement := range program.Statements {
		result = i.evalStatement(statement)
	}

	return result
}

func (i *Interpreter) evalStatement(stmt Statement) Value {
	switch s := stmt.(type) {
	case *LetStatement:
		val := i.evalExpression(s.Value)
		i.env.Set(s.Name.Value, val)
		return val
	case *ReturnStatement:
		return i.evalExpression(s.ReturnValue)
	default:
		return Undefined
	}
}

func (i *Interpreter) evalExpression(exp Expression) Value {
	switch e := exp.(type) {
	case *NumberLiteral:
		return Value{Type: TypeNumber, Data: e.Value}
	case *StringLiteral:
		return Value{Type: TypeString, Data: e.Value}
	case *BooleanLiteral:
		return Value{Type: TypeBoolean, value: e.Value}
	case *Identifier:
		if val, ok := i.env.Get(e.Value); ok {
			return val
		}
		return Undefined
	case *PrefixExpression:
		right := i.evalExpression(e.Right)
		switch e.Operator {
		case "!":
			return Value{Type: TypeBoolean, value: !right.ToBoolean()}
		case "-":
			if right.Type == TypeNumber {
				return Value{Type: TypeNumber, Data: -right.Data.(float64)}
			}
		}
		return Undefined
	case *InfixExpression:
		left := i.evalExpression(e.Left)
		right := i.evalExpression(e.Right)

		switch e.Operator {
		case "+":
			return left.Add(right)
		case "-":
			return left.Subtract(right)
		case "*":
			return left.Multiply(right)
		case "/":
			return left.Divide(right)
		}
		return Undefined
	case *FunctionLiteral:
		return Value{
			Type: TypeFunction,
			Data: &Function{
				Parameters: e.Parameters,
				Body:       e.Body,
				Env:        i.env.store,
			},
		}
	case *CallExpression:
		fn := i.evalExpression(e.Function)
		args := make([]Value, len(e.Arguments))
		for idx, arg := range e.Arguments {
			args[idx] = i.evalExpression(arg)
		}
		return fn.Call(args...)
	default:
		return Undefined
	}
}
