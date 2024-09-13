package main

import (
	"os"

	"monkey/interpreter"
)

func main() {
	if err := interpreter.Interpreter(os.Args[1]); err != nil {
		panic(err)
	}
}
