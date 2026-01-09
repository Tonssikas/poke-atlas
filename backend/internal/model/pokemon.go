package model

type Pokemon struct {
	ID                     string               `json:"id"`
	Name                   string               `json:"name"`
	BaseExperience         int                  `json:"base_experience"`
	Height                 int                  `json:"height"`
	IsDefault              bool                 `json:"is_default"`
	Order                  int                  `json:"order"`
	Weight                 int                  `json:"weight"`
	Abilities              []PokemonAbility     `json:"abilities"`
	Forms                  []NamedResource      `json:"forms"`
	GameIndices            []VersionGameIndex   `json:"game_indices"`
	HeldItems              []HeldItem           `json:"held_items"`
	LocationAreaEncounters string               `json:"location_area_encounters"`
	Moves                  []PokemonMove        `json:"moves"`
	PastTypes              []PokemonTypePast    `json:"past_types"`
	PastAbilities          []PokemonAbilityPast `json:"past_abilities"`
	Sprites                PokemonSprites       `json:"sprites"`
	Cries                  PokemonCries         `json:"cries"`
	Species                NamedResource        `json:"species"`
	Stats                  []PokemonStat        `json:"stats"`
	Types                  []PokemonType        `json:"types"`
}

type NamedResource struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type PokemonAbility struct {
	Slot      int           `json:"slot"`
	Is_hidden bool          `json:"is_hidden"`
	Ability   NamedResource `json:"ability"`
}

type VersionGameIndex struct {
	Game_index uint16        `json:"game_index"`
	Version    NamedResource `json:"version"`
}

type HeldItem struct {
	Item           NamedResource          `json:"item"`
	VersionDetails PokemonHeldItemVersion `json:"version_details"`
}

type Item struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type PokemonHeldItemVersion struct {
	Rarity  int           `json:"rarity"`
	Version NamedResource `json:"version"`
}

type PokemonMove struct {
	Move                NamedResource      `json:"move"`
	VersionGroupDetails PokemonMoveVersion `json:"version_group_details"`
}

type PokemonMoveVersion struct {
	LevelLearnedAt  uint16        `json:"level_learned_at"`
	VersionGroup    NamedResource `json:"version_group"`
	MoveLearnMethod NamedResource `json:"move_learn_method"`
	Order           int           `json:"order"`
}

type PokemonTypePast struct {
	Generation NamedResource `json:"generation"`
	Types      []PokemonType `json:"types"`
}

type PokemonType struct {
	Slot int           `json:"integer"`
	Type NamedResource `json:"type"`
}

type PokemonAbilityPast struct {
	Generation NamedResource    `json:"generation"`
	Abilities  []PokemonAbility `json:"abilities"`
}

type PokemonSprites struct {
	FrontDefault     string `json:"fromt_default"`
	FrontShiny       string `json:"front_shiny"`
	FrontFemale      string `json:"front_female"`
	FrontFemaleShiny string `json:"front_female_shiny"`
	BackDefault      string `json:"back_default"`
	BackShiny        string `json:"back_shiny"`
	BackFemale       string `json:"back_female"`
	BackShinyFemale  string `json:"back_shiny_female"`
}

type PokemonCries struct {
	Latest string `json:"latest"`
	Legacy string `json:"legacy"`
}

type PokemonStat struct {
	Stat     NamedResource `json:"stat"`
	Effort   int           `json:"effort"`
	BaseStat int           `json:"base_stat"`
}
