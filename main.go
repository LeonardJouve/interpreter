package main

import (
	"leonardjouve/repl"
	"os"
)

func main() {
	repl.Start(os.Stdin, os.Stdout)
}
