package runtime

import (
	"mini-js/engine"
	"time"
)

func (r *Runtime) setTimeout(args ...engine.Value) engine.Value {
	if len(args) < 2 {
		return engine.Undefined
	}

	callback := args[0]
	delay := args[1].ToNumber()

	if !callback.IsFunction() {
		return engine.Undefined
	}

	r.eventLoop.AddTask(func() {
		callback.Call()
	}, time.Duration(delay)*time.Millisecond)

	return engine.Undefined
}
