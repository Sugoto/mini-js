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
	return Value{Type: TypeUndefined}, nil
}
