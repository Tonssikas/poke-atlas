package repository

import (
	"context"
	"fmt"
	"log"
	"poke-atlas/web-service/internal/model"
	"poke-atlas/web-service/internal/pokeapi"
	"poke-atlas/web-service/internal/store"
)

type Repository interface {
	GetPokemon(ctx context.Context, name string) (model.Pokemon_summary, error)
}

type repository struct {
	pokeAPIClient pokeapi.PokeAPIClient
	database      store.Database
}

func NewRepository(pokeAPIClient pokeapi.PokeAPIClient, db store.Database) Repository {
	newRepository := &repository{
		pokeAPIClient: pokeAPIClient,
		database:      db,
	}

	return newRepository
}

func (r *repository) GetPokemon(ctx context.Context, name string) (model.Pokemon_summary, error) {

	// Check database first
	pokemon, err := r.database.GetPokemon(ctx, name)
	if err == nil {
		log.Printf("Pokemon %s found in the database\n", name)
		return pokemon, nil
	}

	// If not found in local database, fetch from PokeAPI
	response, err := r.fetchFromPokeAPI(ctx, name)
	if err != nil {
		return model.Pokemon_summary{}, err
	}

	pokemon = model.Pokemon_summary{
		ID:     response.ID,
		Name:   response.Name,
		Weight: response.Weight,
		Height: response.Height,
	}

	return pokemon, nil
}

func (r *repository) fetchFromPokeAPI(ctx context.Context, name string) (model.Pokemon, error) {

	log.Println(fmt.Println("Fetching from api..."))
	pokemon, err := r.pokeAPIClient.GetPokemon(ctx, name)

	if err != nil {
		return model.Pokemon{}, err
	}

	// If pokemon was not in database, but found in pokeAPI add to database
	err = r.database.AddPokemon(ctx, pokemon)

	if err != nil {
		log.Println(err.Error())
	}

	return pokemon, nil
}
