import type { PageServerLoad } from "./$types";
import { error } from '@sveltejs/kit'
import type { PokemonDetailed } from '$lib/model/PokemonDetailed.js'

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

    const response = await fetch(`http://backend:8080/pokemondetailed/${parsedInt}`);
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