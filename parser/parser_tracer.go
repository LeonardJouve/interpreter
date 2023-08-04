package parser

import (
	"fmt"
	"strings"
)

var nestLevel int = 0

const INDENT string = "\t"

func getIndent() string {
	return strings.Repeat(INDENT, nestLevel-1)
}

func printTrace(msg string) {
	fmt.Printf("%s%s\n", getIndent(), msg)
}

func trace(msg string) string {
	nestLevel += 1
	printTrace("BEGIN " + msg)
	return msg
}

func untrace(msg string) {
	printTrace("END " + msg)
	nestLevel -= 1
}
