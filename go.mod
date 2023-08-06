module github.com/LeonardJouve/interpreter

replace leonardjouve/token => ./token

replace leonardjouve/lexer => ./lexer

replace leonardjouve/repl => ./repl

replace leonardjouve/evaluator => ./evaluator

replace leonardjouve/object => ./object

replace leonardjouve/parser => ./parser

replace leonardjouve/ast => ./ast

go 1.20

require leonardjouve/repl v0.0.0-00010101000000-000000000000

require (
	leonardjouve/ast v0.0.0-00010101000000-000000000000 // indirect
	leonardjouve/evaluator v0.0.0-00010101000000-000000000000 // indirect
	leonardjouve/lexer v0.0.0-00010101000000-000000000000 // indirect
	leonardjouve/object v0.0.0-00010101000000-000000000000 // indirect
	leonardjouve/parser v0.0.0-00010101000000-000000000000 // indirect
	leonardjouve/token v0.0.0-00010101000000-000000000000 // indirect
)
