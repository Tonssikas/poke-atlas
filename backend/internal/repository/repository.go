package repository

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"poke-atlas/web-service/internal/model"
	"poke-atlas/web-service/internal/pokeapi"
	"poke-atlas/web-service/internal/store"
	"strconv"
)

type Repository interface {
	GetPokemon(ctx context.Context, name string) (model.Pokemon_summary, error)
	GetPokemons(ctx context.Context, offset int, limit int) ([]model.Pokemon_summary, error)
	GetPokemonDetailed(ctx context.Context, id int) (model.Pokemon_details, error)
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

	log.Print(pokemon)

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
		ID:        response.ID,
		Name:      response.Name,
		Weight:    response.Weight,
		Height:    response.Height,
		SpriteUrl: response.Sprites.FrontDefault,
	}

	return pokemon, nil
}

func (r *repository) GetPokemons(ctx context.Context, offset int, limit int) ([]model.Pokemon_summary, error) {

	// Check database first
	pokemons, err := r.database.GetPokemons(ctx, offset, limit)
	if (err == nil) && (len(pokemons) == limit) {
		log.Print(len(pokemons))
		log.Println("Pokemons found in database")
		return pokemons, nil
	}

	// Fetch from pokeapi
	log.Println("Fetching from api...")
	response, err := r.pokeAPIClient.GetPokemons(ctx, offset, limit)
	if err != nil {
		log.Println("Failed to fetch pokemons from api", err.Error())
		return nil, err
	}

	for _, pokemon := range response {
		//log.Printf("Adding pokemon %s to database", pokemon.Name)
		err = r.database.AddPokemon(ctx, pokemon)
		if err != nil {
			log.Println("Failed to insert pokemon to db", err.Error())
			return nil, err
		}
	}

	// Fetch from database again after insert
	pokemons, err = r.database.GetPokemons(ctx, offset, limit)
	if (err == nil) && (len(pokemons) > 0) {
		log.Println("Pokemons found in database")
		return pokemons, nil
	}

	/*
		pokemons = make([]model.Pokemon_summary, len(response))
		for i, p := range response {
			pokemons[i] = model.Pokemon_summary{
				ID:     p.ID,
				Name:   p.Name,
				Weight: p.Weight,
				Height: p.Height,
			}
		}
		// Sort pokemons fetched from API by ID
		sort.Slice(pokemons, func(i int, j int) bool {
			return pokemons[i].ID < pokemons[j].ID
		})*/

	return pokemons, nil
}

func (r *repository) GetPokemonDetailed(ctx context.Context, id int) (model.Pokemon_details, error) {
	// Database
	pokemon, err := r.database.GetPokemonDetailed(ctx, id)

	//log.Printf("GetPokemonDetailed err: %v, type: %T", err, err)
	//log.Printf("Is sql.ErrNoRows? %v", errors.Is(err, sql.ErrNoRows))

	if errors.Is(err, sql.ErrNoRows) {
		log.Print("pokemon not found in the database!")

		// We can convert to string because pokeAPI supports querying both name and id
		fetchedPokemon, err := r.pokeAPIClient.GetPokemon(ctx, strconv.Itoa(id))
		if err != nil {
			return model.Pokemon_details{}, err
		}

		err = r.database.AddPokemon(ctx, fetchedPokemon)
		if err != nil {
			return model.Pokemon_details{}, err
		}

		// Fetch the detailed view from database after adding
		pokemon, err = r.database.GetPokemonDetailed(ctx, id)
	}

	// Check if we need to fetch evolution chain
	// We need it if we don't have ANY evolution data for this pokemon
	needsEvolutionChain := len(pokemon.EvolutionChain) == 0

	if needsEvolutionChain {
		log.Println("fetching evolution chain from pokeapi...")
		evoChain, err := r.pokeAPIClient.GetEvolutionChain(ctx, id)

		if err != nil {
			// If evolution chain fetch fails, return pokemon without it
			// rather than failing the entire request
			log.Printf("Failed to fetch evolution chain for pokemon %d: %v", id, err)
			return pokemon, nil
		}

		// Add evolution chain data to db
		err = r.database.AddEvolutionChain(ctx, evoChain)

		if err != nil {
			// Fail here most likely indicates that not all pokemon data exists in the database --> FOREIGN KEY contraint failed

			log.Printf("Failed to add evolution chain to db: %v", err)

			// In case of error try fetching possible missing pokemons
			log.Printf("Attempting to fetch missing pokemons...")

			missing := extractNamessFromEvolutionChain(evoChain)

			for _, name := range missing {
				// Fetch missing pokemons
				fetchedPokemon, err := r.pokeAPIClient.GetPokemon(ctx, name)
				if err != nil {
					log.Printf("Failed to fetch missing pokemons: %v", err)
					break
				}
				// Add missing pokemons to db
				err = r.database.AddPokemon(ctx, fetchedPokemon)
				log.Println("adding pokemon: ", name)
				if err != nil {
					log.Printf("Failed to add pokemon to database: %v", err)
				}
			}

			// Retry adding evolution chain to db
			log.Print("Retrying adding evolution chain to db...")
			err = r.database.AddEvolutionChain(ctx, evoChain)
			if err != nil {
				log.Printf("Failed to add evolution chain to db: %v", err)
			}
		}

		// Fetch again to get the evolution data
		pokemon, err = r.database.GetPokemonDetailed(ctx, id)

		if err != nil {
			return model.Pokemon_details{}, err
		}
	}

	return pokemon, nil
}

func extractNamessFromEvolutionChain(chain model.Evolution_chain) []string {
	ids := make(map[string]bool)

	var traverse func(link model.ChainLink)
	traverse = func(link model.ChainLink) {
		ids[link.Species.Name] = true
		for _, evo := range link.EvolvesTo {
			traverse(evo)
		}
	}

	traverse(chain.Chain)

	result := make([]string, 0, len(ids))
	for id := range ids {
		result = append(result, id)
	}
	return result
}
