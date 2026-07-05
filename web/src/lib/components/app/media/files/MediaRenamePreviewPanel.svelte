<script lang="ts">
	import RefreshCwIcon from '@lucide/svelte/icons/refresh-cw';
	import { Button } from '$lib/components/ui/button';
	import { Badge } from '$lib/components/ui/badge';
	import type { MediaRenamePreviewRow } from '$lib/settings/types';

	interface Props {
		rows: MediaRenamePreviewRow[];
		loading: boolean;
		errorMessage?: string;
		onPreview: () => void;
	}

	let { rows, loading, errorMessage, onPreview }: Props = $props();

	function statusVariant(status: MediaRenamePreviewRow['status']) {
		if (status === 'safe') return 'secondary';
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
		<Button type="button" variant="outline" onclick={onPreview} disabled={loading}>
			<RefreshCwIcon aria-hidden="true" />
			{loading ? 'Previewing' : 'Preview'}
		</Button>
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
