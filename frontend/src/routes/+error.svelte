<script lang="ts">
	import { page } from '$app/state';
</script>

<h1>Error {page.status}: {page.error?.message}</h1>

<div class="justify-items-center">
	<img src="/src/lib/images/pokeball.png" alt="pokeball" class="h-auto max-h-50 w-auto max-w-50" />
</div>
