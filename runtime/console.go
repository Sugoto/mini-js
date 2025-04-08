package runtime

import (
	"fmt"
	"mini-js/engine"
)

func ConsoleLog(args ...engine.Value) engine.Value {
	for _, arg := range args {
		fmt.Print(arg.ToString(), " ")
	}
	fmt.Println()
	return engine.Undefined
}
