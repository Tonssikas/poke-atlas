<script lang="ts">
	import github from '$lib/images/github.svg';
	import pokeball from '$lib/images/pokeball.png';
</script>

<header>
	<div class="corner">
		<a href="https://svelte.dev/docs/kit">
			<img src={pokeball} alt="Pokeball" />
		</a>
	</div>

	

	<div class="corner">
		<a href="https://github.com/Tonssikas/poke-atlas">
			<img src={github} alt="GitHub" />
		</a>
	</div>
</header>

<style>
	header {
		display: flex;
		justify-content: space-between;
	}

	.corner {
		width: 3em;
		height: 3em;
	}

	.corner a {
		display: flex;
		align-items: center;
		justify-content: center;
		width: 100%;
		height: 100%;
	}

	.corner img {
		width: 2em;
		height: 2em;
		object-fit: contain;
	}

	a:hover {
		color: var(--color-theme-1);
	}
</style>
