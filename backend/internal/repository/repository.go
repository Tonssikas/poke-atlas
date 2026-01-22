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
	GetPokemons(ctx context.Context, offset int) ([]model.Pokemon_summary, error)
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

	log.Println("Fetching from api...")
	response, err := r.pokeAPIClient.GetPokemon(ctx, name)
	if err != nil {
		log.Println("Failed to fetch pokemon from api", err.Error())
		return model.Pokemon_summary{}, err
	}

	// If pokemon was not in database, but found in pokeAPI add to database
	err = r.database.AddPokemon(ctx, response)

	if err != nil {
		log.Println("Failed to insert pokemon to db", err.Error())
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

func (r *repository) GetPokemons(ctx context.Context, offset int) ([]model.Pokemon_summary, error) {

	// Check database first
	pokemons, err := r.database.GetPokemons(ctx, offset)
	if (err == nil) && (len(pokemons) > 0) && (pokemons[len(pokemons)-1].ID == offset+20) {
		log.Println("Pokemons found in database")
		return pokemons, nil
	}

	// Fetch from pokeapi
	log.Println(fmt.Println("Fetching from api..."))
	response, err := r.pokeAPIClient.GetPokemons(ctx, offset)
	if err != nil {
		log.Println("Failed to fetch pokemons from api", err.Error())
		return nil, err
	}

	for _, pokemon := range response {
		err = r.database.AddPokemon(ctx, pokemon)
		if err != nil {
			log.Println("Failed to insert pokemon to db", err.Error())
			return nil, err
		}
	}

	pokemons = make([]model.Pokemon_summary, len(response))
	for i, p := range response {
		pokemons[i] = model.Pokemon_summary{
			ID:     p.ID,
			Name:   p.Name,
			Weight: p.Weight,
			Height: p.Height,
		}
	}

	return pokemons, nil
}
