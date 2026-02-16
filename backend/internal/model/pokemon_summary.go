package model

type Pokemon_summary struct {
	ID        int      `json:"id"`
	Name      string   `json:"name"`
	Weight    int      `json:"weight"`
	Height    int      `json:"height"`
	SpriteUrl string   `json:"sprite_url"`
	Types     []string `json:"types"`
}
