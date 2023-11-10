package id

import (
	"github.com/google/uuid"
)

type Session struct{ BaseID }
type Room struct{ BaseID }
type Game struct{ BaseID }

type BaseID struct {
	uuid.UUID
}

type ID interface {
	Session | Room | Game
}

func New[T ID]() T {
	return T{BaseID{uuid.New()}}
}
