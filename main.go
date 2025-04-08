package main

import (
	"mini-js/runtime"
)

func main() {
	rt := runtime.NewRuntime()
	defer rt.Close()

	// rt.EnableDebug()

	_, err := rt.Execute(`
		let add = function(a, b) {
			return a + b;
		};
		
		let result = add(10, 20);
		console.log("The result is:", result);
	`)

	if err != nil {
		panic(err)
	}
}
