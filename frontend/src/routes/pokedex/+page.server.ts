import type { PokemonSummary } from '$lib/model/PokemonSummary'

export async function load({ fetch }) {
    // Fetch only 1 hard coded pokemon for now
    const response = await fetch ('http://localhost:8080/pokemon/bulbasaur');
    const pokemon: PokemonSummary = await response.json();

    console.log(pokemon);
    
    return {
        pokemon: [pokemon] // Wrap in array for the {#each} in +page.svelte
    }
}