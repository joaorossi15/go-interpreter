package repl

import (
	"bufio"
	"fmt"
	"io"

	"monkey/eval"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
)

const MONKEY_FACE = ` 
     w  c(..)o   (
      \__(-)    __)
          /\   (
         /(_)___)
         w /|
          | \
         m  m
`

func REPL(in io.Reader, out io.Writer) error {
	scanner := bufio.NewScanner(in)
	env := object.NewEnviroment()

	for {
		fmt.Printf(">> ")

		scanned := scanner.Scan()

		if !scanned {
			if err := scanner.Err(); err != nil {
				return err
			}
			return nil
		}

		line := scanner.Text()

		if line == "exit" {
			return nil
		}

		l := lexer.New(line)
		p := parser.NewParser(l)
		prog := p.ParseProgram()

		if len(p.Errors()) != 0 {
			printErrors(out, p.Errors())
			continue
		}

		/*
			for _, s := range prog.Statements {
				fmt.Print(s.String())
			}
		*/

		evaluated := eval.Eval(prog, env)

		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}

	}
}

func printErrors(out io.Writer, errors []string) {
	io.WriteString(out, MONKEY_FACE)
	io.WriteString(out, "Looks like we ran into some monkey business here...\nparser errors:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
