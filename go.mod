module github.com/leonardjouve/interpreter

replace leonardjouve/token => ./token

replace leonardjouve/lexer => ./lexer

replace leonardjouve/repl => ./repl

go 1.20

require leonardjouve/repl v0.0.0-00010101000000-000000000000

require (
	leonardjouve/lexer v0.0.0-00010101000000-000000000000 // indirect
	leonardjouve/token v0.0.0-00010101000000-000000000000 // indirect
)
