package pokeapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"poke-atlas/web-service/internal/model"
	"sync"
)

type PokeAPIClient interface {
	GetPokemon(ctx context.Context, name string) (model.Pokemon, error)
	GetPokemons(ctx context.Context, offset int) ([]model.Pokemon, error)
}

type pokeAPIClient struct {
	client *http.Client
}

func NewPokeAPIClient(httpClient *http.Client) PokeAPIClient {
	newPokeAPIClient := &pokeAPIClient{
		client: httpClient,
	}

	return newPokeAPIClient
}

func (c *pokeAPIClient) GetPokemon(ctx context.Context, name string) (model.Pokemon, error) {
	url := fmt.Sprintf("http://pokeapi.co/api/v2/pokemon/%s", name)
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)

	if err != nil {
		return model.Pokemon{}, fmt.Errorf("creating request: %w", err)
	}

	response, err := c.client.Do(request)
	if err != nil {
		return model.Pokemon{}, fmt.Errorf("fetching pokemon: %w", err)
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(response.Body)
		return model.Pokemon{}, fmt.Errorf("PokeAPI returned status %d: %s", response.StatusCode, string(body))
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Print(err.Error())
		return model.Pokemon{}, err
	}

	var pokemon model.Pokemon

	err = json.Unmarshal(body, &pokemon)

	if err != nil {
		fmt.Print(err.Error())
		return model.Pokemon{}, err
	}

	return pokemon, err
}

func (c *pokeAPIClient) GetPokemons(ctx context.Context, offset int) ([]model.Pokemon, error) {
	limit := 20
	// Fetch list of pokemon names by id
	listURL := fmt.Sprintf(
		"https://pokeapi.co/api/v2/pokemon?offset=%d&limit=%d", offset, limit,
	)

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, listURL, nil)
	if err != nil {
		return nil, err
	}

	response, err := c.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Pokeapi returned status %d: %s", response.StatusCode, string(body))
	}

	type PokemonListResponse struct {
		Results []struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"results"`
	}

	var list PokemonListResponse
	if err := json.Unmarshal(body, &list); err != nil {
		return nil, err
	}

	type result struct {
		pokemon model.Pokemon
		err     error
	}

	// Fetch pokemon data concurrently using goroutines
	results := make(chan result, len(list.Results))
	sem := make(chan struct{}, 5) // Limit max concurrent goroutines to 5

	var wg sync.WaitGroup

	for _, entry := range list.Results {
		wg.Add(1)

		go func(url string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
			if err != nil {
				results <- result{err: err}
				return
			}

			resp, err := c.client.Do(req)
			if err != nil {
				results <- result{err: err}
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				results <- result{err: fmt.Errorf("pokemon status %d", resp.StatusCode)}
				return
			}

			var p model.Pokemon
			if err := json.NewDecoder(resp.Body).Decode(&p); err != nil {
				results <- result{err: err}
				return
			}

			results <- result{pokemon: p}
		}(entry.URL)
	}

	wg.Wait()
	close(results)

	// Collect results
	pokemons := make([]model.Pokemon, 0, len(list.Results))
	for res := range results {
		if res.err != nil {
			return nil, res.err
		}
		pokemons = append(pokemons, res.pokemon)
	}

	return pokemons, nil
}
