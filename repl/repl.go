package repl

import (
	"bufio"
	"fmt"
	"io"
	"leonardjouve/evaluator"
	"leonardjouve/lexer"
	"leonardjouve/object"
	"leonardjouve/parser"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironement()
	macroEnv := object.NewEnvironement()

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

		evaluator.DefineMacros(program, macroEnv)
		expanded := evaluator.ExpandMacros(program, macroEnv)

		eval := evaluator.Eval(expanded, env)
		if eval == nil {
			continue
		}
		io.WriteString(out, eval.Inspect()+"\n")
	}
}

func printParserErrors(out io.Writer, errors []string) {
	for _, err := range errors {
		io.WriteString(out, "\t"+err+"\n")
	}
}
