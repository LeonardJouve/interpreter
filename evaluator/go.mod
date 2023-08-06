module evaluator

replace leonardjouve/ast => ../ast

replace leonardjouve/object => ../object

replace leonardjouve/token => ../token

replace leonardjouve/lexer => ../lexer

replace leonardjouve/parser => ../parser

go 1.20

require (
	leonardjouve/ast v0.0.0-00010101000000-000000000000
	leonardjouve/lexer v0.0.0-00010101000000-000000000000
	leonardjouve/object v0.0.0-00010101000000-000000000000
	leonardjouve/parser v0.0.0-00010101000000-000000000000
)

require leonardjouve/token v0.0.0-00010101000000-000000000000 // indirect
