<script lang="ts">
	import TypeIcon from '$lib/components/TypeIcon.svelte';
	import type { evolutionChain } from '$lib/model/PokemonDetailed.js';
	import type { PageData } from './$types';

	let { data } = $props();

	type EvolutionNode = {
		id: number;
		name: string;
		min_level?: number;
		trigger?: string;
		branches?: EvolutionNode[];
	};

	function getEvolutionChain(): EvolutionNode[] {
		if (!data.pokemon?.evolution_chain || data.pokemon?.evolution_chain.length === 0) {
			return [
				{
					id: data.pokemon?.id ?? 0,
					name: data.pokemon?.name ?? ''
				}
			];
		}

		const evolutionMap: Map<number, evolutionChain[]> = new Map();

		data.pokemon.evolution_chain.forEach((evo) => {
			if (!evolutionMap.has(evo.pokemon_id)) {
				evolutionMap.set(evo.pokemon_id, []);
			}
			evolutionMap.get(evo.pokemon_id)!.push(evo);
		});

		// Find the root pokemon (First pokemon of the chain)
		// always first entry of the map?
		const allEvolvesToIds = new Set(data.pokemon.evolution_chain.map((e) => e.evolves_to_id));

		const root = data.pokemon.evolution_chain.find((e) => !allEvolvesToIds.has(e.pokemon_id));

		if (!root) return [];

		function createNode(pokemonId: number, pokemonName: string): EvolutionNode {
			const node: EvolutionNode = {
				id: pokemonId,
				name: pokemonName
			};

			// Get all evolutions from this pokemon (handles branches)
			const evolutions = evolutionMap.get(pokemonId);

			if (evolutions && evolutions.length > 0) {
				node.branches = evolutions.map((evo) => ({
					...createNode(evo.evolves_to_id, evo.evolves_to_name),
					min_level: evo.min_level,
					trigger: evo.trigger_name
				}));
			}

			return node;
		}
		return [createNode(root.pokemon_id, root.pokemon_name)];
	}
</script>

<svelte:head>
	<title>Pokemon-details</title>
	<meta name="description" content="Pokemon details page" />
	<meta name="viewport" content="width=device-width, initial-scale=1.0" />
</svelte:head>

<h1 class="capitalize">{data.pokemon?.name}</h1>

<section>
	<div class="flex w-5/6 flex-col-reverse justify-between lg:flex-row">
		<div class="self-center">
			<div class="rounded-md bg-surface-200 p-5 md:mb-3 lg:mb-10">
				<p>Pokemon ID: {data.pokemon?.id}</p>
				<p class="capitalize">Name: {data.pokemon?.name}</p>
				<p>Weight: {data.pokemon?.weight}</p>
				<p>Height: {data.pokemon?.height}</p>
			</div>
			<div class="md: mb-10 rounded-md bg-surface-200 p-5">
				{#each data.pokemon?.stats as stat}
					<p class="capitalize">{stat.stat_name}: {stat.base_stat}</p>
				{/each}
			</div>
			<div class="lg: mb-5 self-center rounded-md bg-surface-200 p-5 md:mb-10">
				<p class="text-center text-xl">Types</p>
				<div class="flex flex-row justify-center gap-x-10">
					{#each data.pokemon?.types as type}
						<div class="self-center">
							<TypeIcon {type} classString="size-20"></TypeIcon>
						</div>
					{/each}
				</div>
			</div>
		</div>
		<div class="self-center">
			<img
				src={data.pokemon?.sprite_url}
				alt={data.pokemon?.name}
				class="h-65 w-auto md:h-100 lg:h-150"
			/>
		</div>
	</div>

	<div class="w-5/6 lg:flex lg:justify-between lg:gap-x-5">
		<div class="justify-items-center rounded-md bg-surface-200 p-5 lg:flex-1">
			<h4 class="mb-4">Evolution chain</h4>
			<div class="flex justify-center">
				{#each getEvolutionChain() as node}
					{@render evolutionNode(node)}
				{/each}
			</div>
		</div>

		<div class="mt-10 justify-items-center rounded-md bg-surface-200 p-5 lg:mt-0 lg:flex-1">
			<h4>Weaknesses and strengths</h4>
		</div>
	</div>
</section>

{#snippet evolutionNode(node: EvolutionNode)}
	<div class="flex items-center gap-4">
		<div class="flex flex-col items-center">
			<a href="/pokedex/{node.id}" class="transition-opacity hover:opacity-80">
				<div
					class="rounded-lg border-2 border-transparent bg-surface-100 p-3 hover:border-primary-500"
				>
					<img
						src="https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/{node.id}.png"
						alt={node.name}
					/>
					<p class="text-center font-semibold capitalize">{node.name}</p>
					<p class="text-center text-xs text-gray-500">#{node.id}</p>
				</div>
			</a>
		</div>

		{#if node.branches && node.branches.length > 0}
			<div class="flex flex-col gap-3">
				{#each node.branches as branch}
					<div class="flex items-center gap-2">
						<div class="flex flex-col items-center px-2 text-sm">
							<span class="text-xl">→</span>
							{#if branch.min_level}
								<span class="text-xs text-gray-600">Lv. {branch.min_level}</span>
							{:else if branch.trigger}
								<span class="text-xs text-gray-600 capitalize"
									>{branch.trigger.replace('-', ' ')}</span
								>
							{/if}
						</div>
						{@render evolutionNode(branch)}
					</div>
				{/each}
			</div>
		{/if}
	</div>
{/snippet}

<style>
	section {
		display: flex;
		flex-direction: column;
		justify-content: center;
		align-items: center;
		flex: 0.6;
	}
</style>
