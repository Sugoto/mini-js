package runtime

import (
	"mini-js/engine"
)

type Runtime struct {
	interpreter *engine.Interpreter
	eventLoop   *EventLoop
}

func NewRuntime() *Runtime {
	r := &Runtime{
		interpreter: engine.NewInterpreter(),
		eventLoop:   NewEventLoop(),
	}

	r.injectGlobals()
	return r
}

func (r *Runtime) injectGlobals() {
	// Inject console.log
	r.interpreter.SetGlobal("console", map[string]interface{}{
		"log": ConsoleLog,
	})

	// Inject setTimeout
	r.interpreter.SetGlobal("setTimeout", r.setTimeout)
}

func (r *Runtime) Execute(code string) (engine.Value, error) {
	return r.interpreter.Eval(code)
}
