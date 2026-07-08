<script lang="ts">
	import { Badge } from '$lib/components/ui/badge';
	import * as Dialog from '$lib/components/ui/dialog';
	import { formatDateTimeWithSeconds } from '$lib/settings/dateFormat';
	import type { SystemJobExecution, SystemJobExecutionLog } from '$lib/settings/types';

	interface Props {
		execution?: SystemJobExecution;
		logs: SystemJobExecutionLog[];
		loading: boolean;
		onClose: () => void;
	}

	let { execution, logs, loading, onClose }: Props = $props();

	function dataText(log: SystemJobExecutionLog) {
		if (!log.data || Object.keys(log.data).length === 0) return '';
		return JSON.stringify(log.data, null, 2);
	}

	function severityClass(severity: string) {
		if (severity === 'error') return 'border-destructive/50 bg-destructive/10 text-destructive';
		if (severity === 'warning') return 'border-amber-500/50 bg-amber-500/10 text-amber-300';
		return 'border-sky-500/50 bg-sky-500/10 text-sky-300';
	}
</script>

<Dialog.Root open={!!execution} onOpenChange={(open) => !open && onClose()}>
	<Dialog.Content
		class="grid max-h-[min(720px,calc(100vh-32px))] w-[min(960px,calc(100vw-32px))] grid-rows-[auto_minmax(0,1fr)] gap-4 sm:max-w-none"
	>
		<Dialog.Header>
			<Dialog.Title>Execution logs</Dialog.Title>
			<Dialog.Description>
				Job {execution?.riverJobId} ({execution?.kind})
			</Dialog.Description>
		</Dialog.Header>

		<div class="min-h-0 overflow-auto rounded-md border border-border">
			{#if loading}
				<p class="m-0 p-4 text-muted-foreground">Loading logs</p>
			{:else if logs.length === 0}
				<p class="m-0 p-4 text-muted-foreground">No logs were recorded for this execution.</p>
			{:else}
				<div class="grid">
					{#each logs as log (log.id)}
						<div class="grid gap-2 border-b border-border p-3 last:border-b-0">
							<div class="flex items-center gap-2">
								<span class="font-mono text-xs text-muted-foreground">
									{formatDateTimeWithSeconds(log.createdAt)}
								</span>
								<Badge variant="outline" class={severityClass(log.severity)}>{log.severity}</Badge>
								<strong class="min-w-0 truncate text-sm">{log.message}</strong>
							</div>
							{#if dataText(log)}
								<pre class="m-0 overflow-auto rounded bg-muted/40 p-2 text-xs leading-5">{dataText(
										log
									)}</pre>
							{/if}
						</div>
					{/each}
				</div>
			{/if}
		</div>
	</Dialog.Content>
</Dialog.Root>
