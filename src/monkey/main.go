package main

import (
	"fmt"
	"os"

	"monkey/repl"
)

func main() {
	fmt.Println("Type any commands")

	if err := repl.REPL(os.Stdin, os.Stdout); err != nil {
		panic(err)
	}
}
