package repl

import (
	"bufio"
	"fmt"
	"io"

	"monkey/lexer"
	"monkey/parser"
	"monkey/token"
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

		for t := l.NextToken(); t.Type != token.EOF; t = l.NextToken() {
			fmt.Printf("%+v\n", t)
		}

		l = lexer.New(line)

		p := parser.NewParser(l)

		prog := p.ParseProgram()

		if len(p.Errors()) != 0 {
			printErrors(out, p.Errors())
			continue
		}

		io.WriteString(out, prog.String())
		io.WriteString(out, "\n")
	}
}

func printErrors(out io.Writer, errors []string) {
	io.WriteString(out, MONKEY_FACE)
	io.WriteString(out, "Looks like we ran into some monkey business here...\nparser errors:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
