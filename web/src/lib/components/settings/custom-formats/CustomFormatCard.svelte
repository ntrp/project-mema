<script lang="ts">
	import CustomFormatSpecChip from '$lib/components/settings/custom-formats/CustomFormatSpecChip.svelte';
	import SettingsRowActionButton from '$lib/components/settings/shared/SettingsRowActionButton.svelte';
	import type { CustomFormat, CustomFormatSpec } from '$lib/settings/types';

	type PreviewSpec = {
		spec: CustomFormatSpec;
		negated: boolean;
	};

	interface Props {
		format: CustomFormat;
		deleting: boolean;
		onEdit: (_format: CustomFormat) => void;
		onDelete: (_id: string) => void | Promise<void>;
	}

	let { format, deleting, onEdit, onDelete }: Props = $props();

	let expanded = $state(false);

	const allSpecs = $derived<PreviewSpec[]>([
		...format.includeSpecs
			.filter((spec) => spec.required)
			.map((spec) => ({ spec, negated: false })),
		...format.excludeSpecs.map((spec) => ({ spec, negated: true })),
		...format.includeSpecs
			.filter((spec) => !spec.required)
			.map((spec) => ({ spec, negated: false }))
	]);
	const previewSpecs = $derived(allSpecs.slice(0, 6));
	const remainingSpecs = $derived(allSpecs.slice(previewSpecs.length));
	const hiddenPreviewCount = $derived(remainingSpecs.length);
	const hiddenSpecIds = $derived(
		remainingSpecs.map((item) => `${format.id}-${item.spec.id}`).join(' ')
	);
</script>

<article class="grid min-h-34 gap-5 rounded-md border border-border bg-muted p-5">
	<div class="flex items-start justify-between gap-3">
		<div class="grid min-w-0">
			<h3 class="m-0 break-words text-2xl font-bold text-muted-foreground">{format.name}</h3>
		</div>
		<div class="flex shrink-0 items-center gap-2">
			<SettingsRowActionButton
				label={`Edit ${format.name}`}
				icon="edit"
				onclick={() => onEdit(format)}
			/>
			<SettingsRowActionButton
				label={`${deleting ? 'Deleting' : 'Delete'} ${format.name}`}
				icon="delete"
				variant="destructive"
				disabled={deleting}
				onclick={() => onDelete(format.id)}
			/>
		</div>
	</div>

	<div class="flex flex-wrap items-start content-start gap-2">
		{#each previewSpecs as item (item.spec.id)}
			<CustomFormatSpecChip spec={item.spec} negated={item.negated} />
		{/each}
		{#if hiddenPreviewCount > 0 && !expanded}
			<button
				type="button"
				class="inline-flex min-h-5 items-start rounded-md bg-secondary px-1.5 text-[11px] font-extrabold text-secondary-foreground hover:bg-secondary/80 focus-visible:ring-ring/50 focus-visible:ring-3 focus-visible:outline-none"
				aria-expanded={expanded}
				aria-controls={hiddenSpecIds || undefined}
				aria-label={`Show ${hiddenPreviewCount} more specs for ${format.name}`}
				onclick={() => (expanded = !expanded)}
			>
				+{hiddenPreviewCount}
			</button>
		{/if}
		{#if expanded}
			{#each remainingSpecs as item (item.spec.id)}
				<span id={`${format.id}-${item.spec.id}`} class="inline-flex items-start">
					<CustomFormatSpecChip spec={item.spec} negated={item.negated} />
				</span>
			{/each}
		{/if}
	</div>
</article>
