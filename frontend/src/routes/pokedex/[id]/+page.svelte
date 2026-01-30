<script lang="ts">
import type { PageData } from './$types';

let { data }: { data: PageData } = $props();
</script>


<h1>Pokemon details page</h1>
<p> Pokemon id: {data.id}</p>