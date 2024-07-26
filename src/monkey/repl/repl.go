package repl

import (
	"bufio"
	"fmt"
	"io"

	"monkey/lexer"
	"monkey/parser"
	"monkey/token"
)

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
		p.PrintParser(prog)
	}
}
