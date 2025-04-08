package engine

type Interpreter struct {
	globals map[string]interface{}
}

func NewInterpreter() *Interpreter {
	return &Interpreter{
		globals: make(map[string]interface{}),
	}
}

func (i *Interpreter) SetGlobal(name string, value interface{}) error {
	i.globals[name] = value
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
			return val.(Value)
		}
		return Undefined
	default:
		return Undefined
	}
}
