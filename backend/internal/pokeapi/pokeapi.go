package pokeapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"poke-atlas/web-service/internal/model"
)

type PokeAPIClient interface {
	GetPokemon(ctx context.Context, name string) (model.Pokemon, error)
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
