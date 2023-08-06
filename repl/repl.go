package repl

import (
	"bufio"
	"fmt"
	"io"
	"leonardjouve/lexer"
	"leonardjouve/parser"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Fprint(out, PROMPT)

		input := scanner.Scan()
		if !input {
			return
		}

		line := scanner.Text()
		lex := lexer.New(line)
		par := parser.New(lex)
		program := par.ParseProgram()

		if len(par.Errors) > 0 {
			printParserErrors(out, par.Errors)
			continue
		}

		io.WriteString(out, "\t"+program.String()+"\n")
	}
}

func printParserErrors(out io.Writer, errors []string) {
	for _, err := range errors {
		io.WriteString(out, "\t"+err+"\n")
	}
}
