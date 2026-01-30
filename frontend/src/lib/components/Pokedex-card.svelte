<script lang="ts">
    import type { PokemonSummary } from '$lib/model/PokemonSummary'
    import TypeIcon from '$lib/components/TypeIcon.svelte'
    let { pokemon, class: className = '' }: { pokemon: PokemonSummary, class?: string } = $props();

    // Local helper function for capitalizing pokemon names
    const capitalize = (str: string) => str.charAt(0).toUpperCase() + str.slice(1);

</script>

<a
href="/pokedex/{pokemon.id}"
class="card rounded-xl w-50 max-w-sm preset-filled-surface-100-900 p-4 text-center justify-items-center {className}">
	<img alt="pokemon" src="/src/lib/images/placeholder-pokemon.png" class="w-30 h-auto"/>
    <p>ID: {pokemon.id}</p>
    <p>{capitalize(pokemon.name)}</p>
    <div class="flex gap-4 pt-3">
        <TypeIcon type="dark" classString="size-8"/>
        <TypeIcon type="water" classString="size-8"/>
    </div>
</a>

<style>
    a {
        text-decoration: none;
        color: initial;
    }
</style>