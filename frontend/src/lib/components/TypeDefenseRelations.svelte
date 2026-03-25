<script lang="ts">

    import { groupEffectiveness } from "$lib/data/TypeEffectiveness";
	import type { PokemonType } from "$lib/icons/pokemon-types";
	import TypeIcon from "./TypeIcon.svelte";

    let { types }: { types: PokemonType[] } = $props();

    
    const effectivenessGroups = $derived(groupEffectiveness(types));
</script>

<div class="mx-auto flex w-fit flex-col gap-5">
    {#if effectivenessGroups.x4.length > 0}
        <div class="grid grid-cols-[4.5rem_auto] items-center gap-x-2">
            <p class="">4x</p>
            <div class="flex flex-wrap items-center gap-2">
                {#each effectivenessGroups.x4 as x4type}
                    <TypeIcon type={x4type} classString="size-8" />
                {/each}
            </div>
        </div>
    {/if}

    {#if effectivenessGroups.x2.length > 0}
        <div class="grid grid-cols-[4.5rem_auto] items-center gap-x-2">
            <p class="">2x</p>
            <div class="flex flex-wrap items-center gap-2">
                {#each effectivenessGroups.x2 as x2type}
                    <TypeIcon type={x2type} classString="size-8" />
                {/each}
            </div>
        </div>
    {/if}

    {#if effectivenessGroups.x05.length > 0}
        <div class="grid grid-cols-[4.5rem_auto] items-center gap-x-2">
            <p class="">0.5x</p>
            <div class="flex flex-wrap items-center gap-2">
                {#each effectivenessGroups.x05 as x05type}
                    <TypeIcon type={x05type} classString="size-8" />
                {/each}
            </div>
        </div>
    {/if}

    {#if effectivenessGroups.x025.length > 0}
        <div class="grid grid-cols-[4.5rem_auto] items-center gap-x-2">
            <p class="">0.25x</p>
            <div class="flex flex-wrap items-center gap-2">
                {#each effectivenessGroups.x025 as x025type}
                    <TypeIcon type={x025type} classString="size-8" />
                {/each}
            </div>
        </div>
    {/if}

    {#if effectivenessGroups.immune.length > 0}
        <div class="grid grid-cols-[4.5rem_auto] items-center gap-x-2">
            <p class="">Immune</p>
            <div class="flex flex-wrap items-center gap-2">
                {#each effectivenessGroups.immune as immune}
                    <TypeIcon type={immune} classString="size-8" />
                {/each}
            </div>
        </div>
    {/if}
</div>