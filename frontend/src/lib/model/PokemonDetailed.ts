import type { PokemonType } from '$lib/icons/pokemon-types';


export type PokemonDetailed = {
    id: number;
    name: string;
    weight: number;
    height: number;
    sprite_url: string;
    types: PokemonType[];
    stats: pokemon_stat[];
    evolution_chain: evolutionChain[];
}

type pokemon_stat = {
    stat_name: statName;
    effort: number;
    base_stat: number;
}

type statName = "attack" | "defense" | "hp" | "special-attack" | "special-defense" | "speed";

export type evolutionChain = {
    pokemon_id: number;
    pokemon_name: string;
    evolves_to_id: number;
    evolves_to_name: string;
    min_level: number;
    trigger_name: string;
}