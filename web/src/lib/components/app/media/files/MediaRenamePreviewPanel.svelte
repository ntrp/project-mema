<script lang="ts">
	import CheckIcon from '@lucide/svelte/icons/check';
	import RefreshCwIcon from '@lucide/svelte/icons/refresh-cw';
	import { Button } from '$lib/components/ui/button';
	import { Badge } from '$lib/components/ui/badge';
	import type { MediaRenamePreviewRow } from '$lib/settings/types';

	interface Props {
		rows: MediaRenamePreviewRow[];
		loading: boolean;
		applying: boolean;
		errorMessage?: string;
		onPreview: () => void;
		onApply: () => void;
	}

	let { rows, loading, applying, errorMessage, onPreview, onApply }: Props = $props();
	const safeCount = $derived(rows.filter((row) => row.status === 'safe').length);

	function statusVariant(status: MediaRenamePreviewRow['status']) {
		if (status === 'safe' || status === 'applied') return 'secondary';
		if (status === 'unchanged') return 'outline';
		return 'destructive';
	}

	function statusLabel(status: MediaRenamePreviewRow['status']) {
		return status === 'missing' ? 'skipped' : status;
	}
</script>

<div class="grid gap-3 rounded-md border bg-card p-4 text-card-foreground">
	<div class="flex flex-wrap items-center justify-between gap-3">
		<div class="grid gap-1">
			<h3 class="m-0 text-base font-semibold">Rename preview</h3>
			<p class="m-0 text-sm text-muted-foreground">
				{rows.length > 0
					? `${rows.length} file${rows.length === 1 ? '' : 's'}`
					: 'No preview loaded'}
			</p>
		</div>
		<div class="flex flex-wrap gap-2">
			<Button type="button" variant="outline" onclick={onPreview} disabled={loading || applying}>
				<RefreshCwIcon aria-hidden="true" />
				{loading ? 'Previewing' : 'Preview'}
			</Button>
			<Button type="button" onclick={onApply} disabled={safeCount === 0 || loading || applying}>
				<CheckIcon aria-hidden="true" />
				{applying ? 'Applying' : 'Apply'}
			</Button>
		</div>
	</div>
	{#if errorMessage}
		<p class="m-0 text-sm text-destructive">{errorMessage}</p>
	{/if}
	{#if rows.length > 0}
		<div class="grid gap-2">
			{#each rows as row (row.currentPath)}
				<div class="grid gap-2 rounded-md border bg-background p-3 md:grid-cols-[8rem_1fr]">
					<Badge variant={statusVariant(row.status)} class="justify-self-start">
						{statusLabel(row.status)}
					</Badge>
					<div class="grid min-w-0 gap-1 text-sm">
						<span class="break-anywhere text-muted-foreground">{row.currentPath}</span>
						<span class="break-anywhere font-medium">{row.proposedPath || '-'}</span>
						{#if row.messages.length > 0}
							<span class="break-anywhere text-muted-foreground">{row.messages.join(' ')}</span>
						{/if}
					</div>
				</div>
			{/each}
		</div>
	{/if}
</div>
