<script lang="ts">
	import { Navigation } from '@skeletonlabs/skeleton-svelte';
	import { MenuIcon, HouseIcon, BookOpenIcon } from '@lucide/svelte';
	import pokeball from '$lib/images/pokeball.png';

	export let layoutRail = true;

	const buttonClasses = 'btn hover:preset-tonal';
	let anchorRail = `${buttonClasses} flex flex-col items-center gap-0.5`;
	let anchorSidebar = `${buttonClasses} justify-start px-2 w-full`;

	function toggleLayout() {
		layoutRail = !layoutRail;
	}
</script>

<Navigation
	layout={layoutRail ? 'rail' : 'sidebar'}
	class="sticky top-0 h-screen border-r border-surface-200-800 transition-[width] duration-200"
>
	<Navigation.Content>
		<Navigation.Header>
			<Navigation.Trigger onclick={toggleLayout}>
				<Navigation.TriggerText>
					<MenuIcon class="size-5 text-surface-900 dark:text-surface-50"></MenuIcon>
				</Navigation.TriggerText>
			</Navigation.Trigger>
		</Navigation.Header>

		<Navigation.Group>
			<Navigation.Menu>
				<Navigation.TriggerAnchor href="/" class={layoutRail ? anchorRail : anchorSidebar}>
					<HouseIcon class="size-5 text-surface-900 dark:text-surface-50"></HouseIcon>
					<Navigation.TriggerText>Home</Navigation.TriggerText>
				</Navigation.TriggerAnchor>
			</Navigation.Menu>
		</Navigation.Group>
		<Navigation.Group>
			<Navigation.Menu>
				<Navigation.TriggerAnchor href="/pokedex" class={layoutRail ? anchorRail : anchorSidebar}>
					<BookOpenIcon class="size-5 text-surface-900 dark:text-surface-50"></BookOpenIcon>
					<Navigation.TriggerText>Pokedex</Navigation.TriggerText>
				</Navigation.TriggerAnchor>
			</Navigation.Menu>
		</Navigation.Group>
	</Navigation.Content>
</Navigation>
