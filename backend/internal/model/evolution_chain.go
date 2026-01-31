package model

type Evolution_chain struct {
	Id              int            `json:"id"`
	BabyTriggerItem *NamedResource `json:"baby_trigger_item"`
	Chain           ChainLink      `json:"chain"`
}

type ChainLink struct {
	IsBaby           bool              `json:"is_baby"`
	Species          NamedResource     `json:"species"`
	EvolutionDetails []EvolutionDetail `json:"evolution_details"`
	EvolvesTo        []ChainLink       `json:"evolves_to"`
}

type EvolutionDetail struct {
	Item                  *NamedResource `json:"item"`
	Trigger               NamedResource  `json:"trigger"`
	Gender                *int           `json:"gender"`
	HeldItem              *NamedResource `json:"held_item"`
	KnownMove             *NamedResource `json:"known_move"`
	KnownMoveType         *NamedResource `json:"known_move_type"`
	Location              *NamedResource `json:"location"`
	MinLevel              *int           `json:"min_level"`
	MinHappiness          *int           `json:"min_happiness"`
	MinBeauty             *int           `json:"min_beauty"`
	MinAffection          *int           `json:"min_affection"`
	NeedsMultiplayer      bool           `json:"needs_multiplayer"`
	NeedsOverworldRain    bool           `json:"needs_overworld_rain"`
	PartySpecies          *NamedResource `json:"party_species"`
	PartyType             *NamedResource `json:"party_type"`
	RelativePhysicalStats *int           `json:"relative_physical_stats"`
	TimeOfDay             string         `json:"time_of_day"`
	TradeSpecies          *NamedResource `json:"trade_species"`
	TurnUpsideDown        bool           `json:"turn_upside_down"`
	UsedMove              *NamedResource `json:"used_move"`
	MinMoveCount          *int           `json:"min_move_count"`
	MinSteps              *int           `json:"min_steps"`
	MinDamageTaken        *int           `json:"min_damage_taken"`
}
