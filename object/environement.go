package object

import "leonardjouve/token"

type Environement struct {
	store map[token.TokenLiteral]Object
	outer *Environement
}

func NewEnvironement() *Environement {
	store := make(map[token.TokenLiteral]Object)
	return &Environement{
		store: store,
		outer: nil,
	}
}

func (env *Environement) Get(identifier token.TokenLiteral) (Object, bool) {
	value, ok := env.store[identifier]
	if !ok && env.outer != nil {
		value, ok = env.outer.Get(identifier)
	}
	return value, ok
}

func (env *Environement) Set(identifier token.TokenLiteral, value Object) {
	env.store[identifier] = value
}

func NewEnclosedEnvironement(outer *Environement) *Environement {
	env := NewEnvironement()
	env.outer = outer
	return env
}
