package model

type Pokemon_details struct {
	ID             int               `json:"id"`
	Name           string            `json:"name"`
	Weight         int               `json:"weight"`
	Height         int               `json:"height"`
	SpriteUrl      string            `json:"sprite_url"`
	Types          []string          `json:"types"`
	Stats          []pokemon_stat    `json:"stats"`
	EvolutionChain []evolution_chain `json:"evolution_chain"`
}

type pokemon_stat struct {
	StatName string `json:"stat_name"`
	Effort   int    `json:"effort"`
	BaseStat int    `json:"base_stat"`
}

type evolution_chain struct {
	PokemonID     int    `json:"pokemon_id"`
	PokemonName   string `json:"pokemon_name"`
	EvolvesToID   int    `json:"evolves_to_id"`
	EvolvesToName string `json:"evolves_to_name"`
	MinLevel      int    `json:"min_level"`
	TriggerName   string `json:"trigger_name"`
}
