import type { PageServerLoad } from "./$types";
import { error } from '@sveltejs/kit';
import type { PokemonDetailed } from '$lib/model/PokemonDetailed.js';
import { API_ADDRESS } from '$env/static/private';

export async function load({ fetch, params, setHeaders }) {
    const parsedInt = parseInt(params.id)
    if (isNaN(parsedInt)) {
        error(400, 'ID must be a valid number');
        return {}
    }

    if (parsedInt < 0) {
        error(400, 'ID must not be negative');
        return {}
    }

    const response = await fetch(`http://${API_ADDRESS}/pokemondetailed/${parsedInt}`);
    const pokemon: PokemonDetailed = await response.json();

    if (pokemon == null) {
        error(404, 'No pokemon found');
    }

    setHeaders({
        "cache-control": "max-age=120" // Cache results for 2 minutes to avoid re-fetching constantly
    })

    console.log(pokemon);

    return {
        pokemon: pokemon
    }
}