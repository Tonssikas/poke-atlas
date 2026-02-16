package store

import (
	"context"
	"poke-atlas/web-service/internal/model"
)

type Database interface {
	InitDB() error
	Close() error
	GetPokemon(ctx context.Context, name string) (model.Pokemon_summary, error)
	GetPokemons(ctx context.Context, offset int, limit int) ([]model.Pokemon_summary, error)
	AddPokemon(ctx context.Context, pokemon model.Pokemon) error
	GetPokemonDetailed(ctx context.Context, id int) (model.Pokemon_details, error)
	AddEvolutionChain(ctx context.Context, chain model.Evolution_chain) error
}
