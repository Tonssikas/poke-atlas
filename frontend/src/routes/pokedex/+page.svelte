<script lang="ts">
	import Card from '$lib/components/Pokedex-card.svelte';
	import { ArrowRightIcon, ArrowLeftIcon } from '@lucide/svelte';
	import { goto, preloadData } from '$app/navigation';
	import { page } from '$app/state';

	let { data } = $props();

	const ITEMS_PER_PAGE = 20;

	const currentOffset = $derived(parseInt(page.url.searchParams.get('offset') || '0'));

	function nextPage() {
		goto(`?offset=${currentOffset + ITEMS_PER_PAGE}`);
	}

	function prevPage() {
		const newOffset = Math.max(0, currentOffset - ITEMS_PER_PAGE);
		goto(`?offset=${newOffset}`);
	}

	// preload next page for smooth user experience
	$effect(() => {
		preloadData(`?offset=${currentOffset + ITEMS_PER_PAGE}`);
	});
</script>

<svelte:head>
	<title>Pokedex-page</title>
	<meta name="description" content="Pokedex" />
</svelte:head>
<h1>Pokedex</h1>

<section>
	<div class="grid gap-6 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-5">
		{#each data.pokemon as pokemon}
			<Card {pokemon} class="pokemon-card border-4 border-transparent"/>
		{/each}
	</div>

	<footer class="flex flex-row gap-8 pt-5">
		{#if currentOffset > 0}
			<button type="button" onclick={prevPage} class="btn-icon preset-filled" title="Next page">
				<ArrowLeftIcon size={24} />
			</button>
		{/if}
		<button type="button" onclick={nextPage} class="btn-icon preset-filled" title="Next page">
			<ArrowRightIcon size={24} />
		</button>
	</footer>
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

	:global(.pokemon-card:hover) {
		background-color: var(--color-surface-400);
		border-style: solid;
		border-color: var(--color-tertiary-100);
	}
</style>
