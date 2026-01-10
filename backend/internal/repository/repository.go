package repository

import (
	"context"
	"poke-atlas/web-service/internal/model"
	"poke-atlas/web-service/internal/pokeapi"
)

type Repository interface {
	GetPokemon(ctx context.Context, name string) (model.Pokemon, error)
}

type repository struct {
	pokeAPIClient pokeapi.PokeAPIClient
}

func NewRepository(pokeAPIClient pokeapi.PokeAPIClient) Repository {
	newRepository := &repository{
		pokeAPIClient: pokeAPIClient,
	}

	return newRepository
}

func (r *repository) GetPokemon(ctx context.Context, name string) (model.Pokemon, error) {

	pokemon, err := r.pokeAPIClient.GetPokemon(ctx, name)

	if err != nil {
		return model.Pokemon{}, err
	}

	return pokemon, nil
}
