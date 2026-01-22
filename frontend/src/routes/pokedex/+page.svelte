<script lang="ts">
    import Card from '$lib/components/Pokedex-card.svelte'

    let { data } = $props();
</script>

<svelte:head>
	<title>Poke-atlas</title>
	<meta name="description" content="Pokedex" />
</svelte:head>

<section>
	<h1>
		Pokedex
	</h1>

	<h2>
		A list of pokemons
	</h2>

    {#each data.pokemon as pokemon}
        <Card {pokemon} />
    {/each}

</section>

<style>
	section {
		display: flex;
		flex-direction: column;
		justify-content: center;
		align-items: center;
		flex: 0.6;
	}

	h1 {
		width: 100%;
	}

</style>
