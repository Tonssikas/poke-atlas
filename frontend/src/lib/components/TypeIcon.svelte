<script lang="ts">
    import { typeIcons, type PokemonType } from '$lib/icons/pokemon-types'

    export let type: PokemonType;
    export let classString: string = '';
</script>

<img
    src={typeIcons[type]}
    alt={type}
    class={ `size-4 ${classString}`}
/>