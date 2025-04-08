package runtime

import (
	"errors"
	"mini-js/engine"
)

type Runtime struct {
	interpreter *engine.Interpreter
	eventLoop   *EventLoop
	isRunning   bool
}

func NewRuntime() *Runtime {
	r := &Runtime{
		interpreter: engine.NewInterpreter(),
		eventLoop:   NewEventLoop(),
		isRunning:   true,
	}

	if err := r.injectGlobals(); err != nil {
		panic("Failed to initialize runtime: " + err.Error())
	}
	return r
}

func (r *Runtime) injectGlobals() error {
	if r.interpreter == nil {
		return errors.New("interpreter not initialized")
	}

	if err := r.interpreter.SetGlobal("console", map[string]interface{}{
		"log": ConsoleLog,
	}); err != nil {
		return err
	}

	if err := r.interpreter.SetGlobal("setTimeout", r.setTimeout); err != nil {
		return err
	}

	return nil
}

func (r *Runtime) Execute(code string) (engine.Value, error) {
	if !r.isRunning {
		return engine.Value{}, errors.New("runtime is stopped")
	}
	if code == "" {
		return engine.Value{}, errors.New("empty code string")
	}
	return r.interpreter.Eval(code)
}

func (r *Runtime) Stop() {
	r.isRunning = false
	r.eventLoop.Clear()
}

func (r *Runtime) Close() error {
	r.Stop()
	return nil
}
