import type { PokemonSummary } from '$lib/model/PokemonSummary'
import { error } from '@sveltejs/kit'
import { API_ADDRESS } from '$env/static/private';

export async function load({ fetch, url, setHeaders }) {
    const offset = url.searchParams.get('offset') || '0';

    const parsedInt = parseInt(offset)
    if (isNaN(parsedInt)) {
        error(400, 'Offset must be a valid number');
        return {}
    }

    if (parsedInt < 0) {
        error(400, 'Offset must not be negative');
        return {}
    }

    const response = await fetch(`http://${API_ADDRESS}/pokemons/${parsedInt}`);
    const pokemon: PokemonSummary[] = await response.json();

    if (pokemon.length === 0) {
        error(404, 'No pokemons found');
    }

    setHeaders({
        "cache-control": "max-age=120" // Cache results for 2 minutes to avoid re-fetching constantly
    })

    console.log(pokemon);

    return {
        pokemon: pokemon
    }
}

