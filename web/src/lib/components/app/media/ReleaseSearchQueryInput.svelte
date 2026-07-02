<script lang="ts">
	import InfoIcon from '@lucide/svelte/icons/info';
	import { Button } from '$lib/components/ui/button';
	import { Checkbox } from '$lib/components/ui/checkbox';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import * as Tooltip from '$lib/components/ui/tooltip';

	interface Props {
		overrideQuery: boolean;
		customQuery: string;
		queryVariants: string[];
		disabled?: boolean;
	}

	let {
		overrideQuery = $bindable(),
		customQuery = $bindable(),
		queryVariants,
		disabled = false
	}: Props = $props();
	let variantsOpen = $state(false);
</script>

<div class="flex min-w-0 flex-wrap items-end gap-3">
	<div class="grid min-w-72 flex-1 gap-2">
		<Label for="release-search-query">Search query</Label>
		<div class="relative">
			<Input
				id="release-search-query"
				class={!overrideQuery ? 'bg-muted pr-10 text-muted-foreground opacity-80' : ''}
				bind:value={customQuery}
				readonly={!overrideQuery}
				{disabled}
				maxlength={500}
			/>
			{#if !overrideQuery}
				<Tooltip.Root bind:open={variantsOpen}>
					<Tooltip.Trigger>
						{#snippet child({ props })}
							<Button
								{...props}
								type="button"
								variant="ghost"
								size="icon-sm"
								class="absolute right-1 top-1/2 -translate-y-1/2 text-muted-foreground"
								aria-label="Show search query variants"
							>
								<InfoIcon aria-hidden="true" />
							</Button>
						{/snippet}
					</Tooltip.Trigger>
					{#if variantsOpen}
						<Tooltip.Content class="max-w-96">
							<div class="grid gap-1 text-left">
								<span class="font-bold">Search branches</span>
								{#each queryVariants as query (query)}
									<span class="font-mono text-xs">{query}</span>
								{/each}
							</div>
						</Tooltip.Content>
					{/if}
				</Tooltip.Root>
			{/if}
		</div>
	</div>
	<div class="flex h-9 items-center gap-2">
		<Checkbox id="release-search-query-override" bind:checked={overrideQuery} {disabled} />
		<Label
			for="release-search-query-override"
			class={disabled ? 'text-muted-foreground opacity-70' : ''}
		>
			Override
		</Label>
	</div>
</div>
