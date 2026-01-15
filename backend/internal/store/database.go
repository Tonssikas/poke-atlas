package store

import (
	"context"
	"poke-atlas/web-service/internal/model"
)

type Database interface {
	InitDB() error
	Close() error
	GetPokemon(ctx context.Context, name string) (model.Pokemon, error)
	AddPokemon(ctx context.Context, pokemon model.Pokemon) error
}
