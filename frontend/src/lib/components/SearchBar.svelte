<script lang="ts">
	import {
		Combobox,
		Portal,
		type ComboboxRootProps,
		useListCollection
	} from '@skeletonlabs/skeleton-svelte';
	import { goto } from '$app/navigation';
	import pokemons from '$lib/data/pokemon-index.json';

	let query = $state('');

	let items = $state(pokemons);

	const collection = $derived(
		useListCollection({
			items: items,
			itemToString: (item) => item.name,
			itemToValue: (item) => item.name
		})
	);

	const onOpenChange = () => {
		items = pokemons;
	};

	const onInputValueChange: ComboboxRootProps['onInputValueChange'] = (event) => {
		const filtered = pokemons.filter((pokemon) =>
			pokemon.name.toLowerCase().includes(event.inputValue.toLowerCase())
		);
		if (filtered.length > 0) {
			items = filtered.slice(0, 10);
		} else {
			items = [];
		}
	};

	const onValueChange: ComboboxRootProps['onValueChange'] = (event) => {
		const selectedPokemon = pokemons.find((pokemon) => pokemon.name === event.value[0]);
		if (selectedPokemon) {
			goto(`/pokedex/${selectedPokemon.id}`);
		}
	};
</script>

	
<Combobox
	class="w-full lg:w-lg"
	placeholder="Search Pokémon..."
	{collection}
	{onOpenChange}
    {onValueChange}
	{onInputValueChange}
	inputBehavior="autohighlight"
>
	<Combobox.Control class="rounded-lg border border-surface-200-800">
		<Combobox.Input
			class="input rounded-lg px-4 py-2 text-surface-900 placeholder:text-surface-500 dark:text-surface-50 dark:placeholder:text-surface-400"
		/>
		<Combobox.Trigger class="btn-icon btn hover:preset-tonal" />
	</Combobox.Control>
	<Portal>
		<Combobox.Positioner>
			<Combobox.Content
				class="max-h-80 overflow-auto card rounded-lg border border-surface-200-800 preset-filled-surface-100-900 shadow-lg"
			>
				{#each items as item (item.name)}
					<Combobox.Item
						{item}
						onclick={() => goto(`/pokedex/${item.id}`)}
						class="cursor-pointer px-4 py-2 text-surface-900 capitalize bg-surface-50 dark:bg-surface-900 dark:text-surface-50
                        data-[highlighted]:!bg-primary-500
                        data-[highlighted]:!text-white
                        hover:preset-bg-primary-500 hover:text-white"
					>
						<Combobox.ItemText>{item.name}</Combobox.ItemText>
						<Combobox.ItemIndicator />
					</Combobox.Item>
				{/each}
			</Combobox.Content>
		</Combobox.Positioner>
	</Portal>
</Combobox>
