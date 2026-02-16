import type { PokemonType } from '$lib/icons/pokemon-types';

export type PokemonSummary = {
    id: number;
    name: string;
    weight: number;
    height: number;
    sprite_url: string;
    types: PokemonType[];
}