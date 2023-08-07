package object

import "leonardjouve/token"

type Environement struct {
	store map[token.TokenLiteral]Object
}

func NewEnvironement() *Environement {
	store := make(map[token.TokenLiteral]Object)
	return &Environement{
		store: store,
	}
}

func (env *Environement) Get(identifier token.TokenLiteral) (Object, bool) {
	value, ok := env.store[identifier]
	return value, ok
}

func (env *Environement) Set(identifier token.TokenLiteral, value Object) {
	env.store[identifier] = value
}
