<script lang="ts">
    import { Progress } from '@skeletonlabs/skeleton-svelte'

    let {statValue, statName}: {statValue: number, statName: string} = $props();
</script>

<div class="flex flex-col items-center w-full">
    <p class="capitalize min-w-fit">{statName}</p>
    <div class="flex items-center gap-2 flex min-w-[250px]">
        <span class="text-sm font-semibold w-6 text-right pr-3">{statValue}</span>
        <Progress value={statValue} max={255} class="flex-1">
            <Progress.Track class="bg-primary-50-950">
                <Progress.Range class="bg-primary-500" />
            </Progress.Track>
        </Progress>
    </div>
</div>


