package engine

import "fmt"

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

func ExtendEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

type Interpreter struct {
	env       *Environment
	debugMode bool
}

func NewInterpreter() *Interpreter {
	return &Interpreter{
		env:       NewEnvironment(),
		debugMode: false,
	}
}

func (i *Interpreter) EnableDebug() {
	i.debugMode = true
}

func (i *Interpreter) DisableDebug() {
	i.debugMode = false
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
		v = Value{Type: TypeBoolean, Data: val}
	default:
		v = Undefined
	}
	i.env.Set(name, v)
	return nil
}

func (i *Interpreter) Eval(code string) (Value, error) {
	if i.debugMode {
		fmt.Println("🔍 Debug: Starting evaluation of code")
	}

	lexer := NewLexer(code)
	if i.debugMode {
		fmt.Println("🔍 Debug: Starting lexical analysis")
	}

	parser := NewParser(lexer)
	if i.debugMode {
		fmt.Println("🔍 Debug: Starting parsing")
	}

	program := parser.ParseProgram()
	if i.debugMode {
		fmt.Println("🔍 Debug: Parsing complete, beginning program evaluation")
	}

	result := i.evalProgram(program)
	if i.debugMode {
		fmt.Printf("🔍 Debug: Program evaluation complete, final result: %v\n", result.ToString())
	}

	return result, nil
}

func (i *Interpreter) evalProgram(program *Program) Value {
	var result Value = Undefined

	for _, statement := range program.Statements {
		if i.debugMode {
			fmt.Printf("🔍 Debug: Evaluating statement: %T\n", statement)
		}
		result = i.evalStatement(statement)
		if result.Type == TypeReturn {
			if i.debugMode {
				fmt.Printf("🔍 Debug: Found return statement, value: %v\n", result.ToString())
			}
			if returnValue, ok := result.Data.(*ReturnValue); ok {
				return returnValue.Value
			}
			return result
		}
	}

	return result
}

func (i *Interpreter) evalStatement(stmt Statement) Value {
	if i.debugMode {
		fmt.Printf("🔍 Debug: Processing statement of type: %T\n", stmt)
	}

	switch s := stmt.(type) {
	case *LetStatement:
		val := i.evalExpression(s.Value)
		if i.debugMode {
			fmt.Printf("🔍 Debug: Let statement - binding '%s' to value: %v\n", s.Name.Value, val.ToString())
		}
		i.env.Set(s.Name.Value, val)
		return Undefined
	case *ReturnStatement:
		val := i.evalExpression(s.ReturnValue)
		if i.debugMode {
			fmt.Printf("🔍 Debug: Return statement - returning value: %v\n", val.ToString())
		}
		return Value{
			Type: TypeReturn,
			Data: &ReturnValue{Value: val},
		}
	case *ExpressionStatement:
		if i.debugMode {
			fmt.Println("🔍 Debug: Evaluating expression statement")
		}
		return i.evalExpression(s.Expression)
	case *BlockStatement:
		if i.debugMode {
			fmt.Println("🔍 Debug: Evaluating block statement")
		}
		return i.evalBlockStatement(s)
	default:
		return Undefined
	}
}

func (i *Interpreter) evalBlockStatement(block *BlockStatement) Value {
	var result Value = Undefined

	for _, statement := range block.Statements {
		result = i.evalStatement(statement)

		if result.Type == TypeReturn {
			return result
		}
	}

	return result
}

func (i *Interpreter) evalExpression(exp Expression) Value {
	if i.debugMode {
		fmt.Printf("🔍 Debug: Evaluating expression of type: %T\n", exp)
	}

	switch e := exp.(type) {
	case *NumberLiteral:
		return Value{Type: TypeNumber, Data: e.Value}
	case *StringLiteral:
		return Value{Type: TypeString, Data: e.Value}
	case *BooleanLiteral:
		return Value{Type: TypeBoolean, Data: e.Value}
	case *Identifier:
		if val, ok := i.env.Get(e.Value); ok {
			return val
		}
		if e.Value == "console" {
			return Value{
				Type: TypeObject,
				Data: "console",
				Properties: map[string]Value{
					"log": Value{
						Type: TypeFunction,
						Data: func(args ...Value) Value {
							for _, arg := range args {
								fmt.Print(arg.ToString(), " ")
							}
							fmt.Println()
							return Undefined
						},
					},
				},
			}
		}
		return Undefined
	case *PrefixExpression:
		right := i.evalExpression(e.Right)
		switch e.Operator {
		case "!":
			return Value{Type: TypeBoolean, Data: !right.ToBoolean()}
		case "-":
			if right.Type == TypeNumber {
				return Value{Type: TypeNumber, Data: -right.Data.(float64)}
			}
		}
		return Undefined
	case *InfixExpression:
		if e.Operator == "." {
			left := i.evalExpression(e.Left)
			if right, ok := e.Right.(*Identifier); ok {
				if left.Properties != nil {
					if prop, ok := left.Properties[right.Value]; ok {
						return prop
					}
				}
				if left.Type == TypeObject && left.Data == "console" && right.Value == "log" {
					return Value{
						Type: TypeFunction,
						Data: func(args ...Value) Value {
							for _, arg := range args {
								fmt.Print(arg.ToString(), " ")
							}
							fmt.Println()
							return Undefined
						},
					}
				}
			}
			return Undefined
		}

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
		params := e.Parameters
		body := e.Body
		return Value{
			Type: TypeFunction,
			Data: &Function{
				Parameters: params,
				Body:       body,
				Env:        i.env,
			},
		}
	case *CallExpression:
		fn := i.evalExpression(e.Function)
		args := make([]Value, len(e.Arguments))
		for idx, arg := range e.Arguments {
			args[idx] = i.evalExpression(arg)
		}
		return i.applyFunction(fn, args)
	default:
		return Undefined
	}
}

func (i *Interpreter) applyFunction(fn Value, args []Value) Value {
	if fn.Type != TypeFunction {
		return Undefined
	}

	switch f := fn.Data.(type) {
	case func(...Value) Value:
		return f(args...)
	case *Function:
		extendedEnv := ExtendEnvironment(f.Env)
		for idx, param := range f.Parameters {
			if idx < len(args) {
				extendedEnv.Set(param.Value, args[idx])
			}
		}
		savedEnv := i.env
		i.env = extendedEnv
		evaluated := i.evalStatement(f.Body)
		i.env = savedEnv
		if evaluated.Type == TypeReturn {
			if returnValue, ok := evaluated.Data.(*ReturnValue); ok {
				return returnValue.Value
			}
		}
		return evaluated
	default:
		return Undefined
	}
}

type ReturnValue struct {
	Value Value
}
