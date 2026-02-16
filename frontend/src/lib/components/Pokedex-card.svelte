<script lang="ts">
    import type { PokemonSummary } from '$lib/model/PokemonSummary'
    import TypeIcon from '$lib/components/TypeIcon.svelte'
    let { pokemon, class: className = '' }: { pokemon: PokemonSummary, class?: string } = $props();

    // Local helper function for capitalizing pokemon names
    const capitalize = (str: string) => str.charAt(0).toUpperCase() + str.slice(1);

</script>

<a
href="/pokedex/{pokemon.id} "
data-sveltekit-preload-data="off"
class="card rounded-xl w-50 max-w-sm preset-filled-surface-100-900 p-4 text-center justify-items-center {className}">
	<img alt="pokemon" src="{pokemon.sprite_url}" class="w-30 h-auto"/>
    <p>ID: {pokemon.id}</p>
    <p>{capitalize(pokemon.name)}</p>
    <div class="flex gap-4 pt-3">
        {#each pokemon.types as type}
            <TypeIcon {type} classString="size-8"/>
        {/each}
    </div>
</a>

<style>
    a {
        text-decoration: none;
        color: initial;
    }
</style>