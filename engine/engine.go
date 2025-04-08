package engine

type Interpreter struct {
	globals map[string]Value
}

func NewInterpreter() *Interpreter {
	return &Interpreter{
		globals: make(map[string]Value),
	}
}

func (i *Interpreter) SetGlobal(name string, value interface{}) error {
	switch v := value.(type) {
	case Value:
		i.globals[name] = v
	case float64:
		i.globals[name] = Value{Type: TypeNumber, Data: v}
	case string:
		i.globals[name] = Value{Type: TypeString, Data: v}
	default:
		i.globals[name] = Undefined
	}
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
		i.globals[s.Name.Value] = val
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
	case *Identifier:
		if val, ok := i.globals[e.Value]; ok {
			return val
		}
		return Undefined
	default:
		return Undefined
	}
}
