package main

import (
	"fmt"
	"mini-js/runtime"
)

func main() {
	rt := runtime.NewRuntime()
	defer rt.Close()

	// Test basic arithmetic
	// result, err := rt.Execute(`
	// 	let x = 5;
	// 	let y = 3;
	// 	x + y;
	// `)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("5 + 3 =", result.ToString())

	// Test function declaration and call
	// result, err := rt.Execute(`
	// 	let add = function(a, b) {
	// 		return a + b;
	// 	};
	// 	add(10, 20);
	// `)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("add(10, 20) =", result.ToString())

	// Test console.log
	_, err := rt.Execute(`
		console.log("Hello from dsfv!");
	`)
	fmt.Println("Done")
	if err != nil {
		panic(err)
	}
}
