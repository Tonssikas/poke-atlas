<script lang="ts">
	import Header from './Header.svelte';
	import Navigation from '../lib/components/Navigation.svelte';
	import './layout.css';
	import { navigating } from '$app/state';

	let { children } = $props();
</script>

<div class="layout">
	<Navigation />

	<div class="main-column">
		<Header />

		{#if navigating.to}
			Navigating to {navigating.to.url.pathname}
		{/if}

		<main>{@render children()}</main>
		<footer>
			<p>Poke-atlas is work in progress</p>
		</footer>
	</div>
</div>

<style>
	.layout {
		display: grid;
		grid-template-columns: auto 1fr;
		min-height: 100vh;
	}

	main {
		flex: 1;
		display: flex;
		flex-direction: column;
		padding: 1rem;
		box-sizing: border-box;
		overflow: auto;
	}

	.main-column {
		display: flex;
		flex-direction: column;
		min-width: 0;
	}

	footer {
		display: flex;
		flex-direction: column;
		justify-content: center;
		align-items: center;
		padding: 12px;
	}

	@media (min-width: 480px) {
		footer {
			padding: 12px 0;
		}
	}
</style>
