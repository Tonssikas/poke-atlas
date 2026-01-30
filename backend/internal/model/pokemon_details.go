package model

type Pokemon_details struct {
	ID     int            `json:"id"`
	Name   string         `json:"name"`
	Weight int            `json:"weight"`
	Height int            `json:"height"`
	Types  []string       `json:"types"`
	Stats  []pokemon_stat `json:"stats"`
}

type pokemon_stat struct {
	Stat_name string `json:"stat_name"`
	Effort    int    `json:"effort"`
	Base_stat int    `json:"base_stat"`
}
