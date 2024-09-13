package main

import (
	"fmt"
	"os"

	"monkey/interpreter"
	"monkey/repl"
)

func main() {
	if len(os.Args) > 1 {
		if os.Args[1] == "i" {
			if len(os.Args) < 3 {
				fmt.Print("File not found...")
				return
			} else if err := interpreter.Interpreter(os.Args[2]); err != nil {
				panic(err)
			}
			return
		} else if os.Args[1] == "r" {
			if err := repl.REPL(os.Stdin, os.Stdout); err != nil {
				panic(err)
			}
			return
		} else {
			fmt.Print("Invalid argument, use i or r")
			return
		}
	}
	fmt.Print("Add i or r as argument")
}
