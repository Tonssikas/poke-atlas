<script lang="ts">
    import type { PokemonSummary } from '$lib/model/PokemonSummary'
    let { pokemon } = $props();
</script>

<a href="/pokedex" 
class="card w-full max-w-sm preset-filled-surface-100-900 p-4 m-1 text-center">
	<p>Pokemon icon here</p>
    <p>ID: {pokemon.id}</p>
    <p>{pokemon.name}</p>
    <p>Type icons</p>
</a>