module leonardjouve/repl

replace leonardjouve/token => ../token

replace leonardjouve/lexer => ../lexer

replace leonardjouve/parser => ../parser

replace leonardjouve/ast => ../ast

replace leonardjouve/evaluator => ../evaluator

replace leonardjouve/object => ../object

go 1.20

require (
	leonardjouve/evaluator v0.0.0-00010101000000-000000000000
	leonardjouve/lexer v0.0.0-00010101000000-000000000000
	leonardjouve/parser v0.0.0-00010101000000-000000000000
)

require (
	leonardjouve/ast v0.0.0-00010101000000-000000000000 // indirect
	leonardjouve/object v0.0.0-00010101000000-000000000000 // indirect
	leonardjouve/token v0.0.0-00010101000000-000000000000 // indirect
)
