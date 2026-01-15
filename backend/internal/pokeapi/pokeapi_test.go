package pokeapi

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"
)

type mockRoundTripper struct {
	fn func(req *http.Request) (*http.Response, error)
}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.fn(req)
}

func TestGetPokemonSuccess(t *testing.T) {
	mockResponse := `{
		"id": 25,
		"name": "bulbasaur"
	}`

	client := &pokeAPIClient{client: &http.Client{
		Transport: &mockRoundTripper{
			fn: func(req *http.Request) (*http.Response, error) {
				if req.Method != http.MethodGet {
					t.Fatalf("Expected GET, got %s", req.Method)
					return nil, nil
				}

				expectedURL := "http://pokeapi.co/api/v2/pokemon/bulbasaur"

				if req.URL.String() != expectedURL {
					t.Fatalf("unexpected URL %s", req.URL)
				}

				// Return a proper mock response
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(mockResponse)),
					Header:     make(http.Header),
				}, nil
			},
		},
	}}
	pokemon, err := client.GetPokemon(context.Background(), "bulbasaur")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if pokemon.Name != "bulbasaur" {
		t.Errorf("expected name 'bulbasaur', got '%s'", pokemon.Name)
	}
}

func TestGetPokemonFail(t *testing.T) {
	mockResponse := `"Not Found"`

	client := &pokeAPIClient{client: &http.Client{
		Transport: &mockRoundTripper{
			fn: func(req *http.Request) (*http.Response, error) {
				if req.Method != http.MethodGet {
					t.Fatalf("Expected GET, got %s", req.Method)
					return nil, nil
				}

				// Return a proper mock response
				return &http.Response{
					StatusCode: http.StatusBadRequest,
					Body:       io.NopCloser(bytes.NewBufferString(mockResponse)),
					Header:     make(http.Header),
				}, nil
			},
		},
	}}
	pokemon, err := client.GetPokemon(context.Background(), "¤")

	if err == nil {
		t.Fatalf("Expected error")
	}

	if pokemon.Name != "" {
		t.Errorf("Expected no pokemon, got %s", pokemon.Name)
	}
}
