<script lang="ts">
	import CirclePlayIcon from '@lucide/svelte/icons/circle-play';
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import * as Table from '$lib/components/ui/table';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import type { ManualFulfillmentAction } from './manualFulfillmentActions';

	interface Props {
		actions: ManualFulfillmentAction[];
		loading?: boolean;
		errorMessage?: string;
	}

	let { actions, loading = false, errorMessage = '' }: Props = $props();

	function disabledReason(action: ManualFulfillmentAction) {
		if (!action.available) return action.blockedReason || 'Unavailable';
		if (action.path.includes('{id}')) return 'Open a matching media or activity item to run';
		return '';
	}
</script>

<div class="overflow-auto rounded-md border border-border">
	<Table.Root class="min-w-full table-auto border-collapse">
		<Table.Header class="bg-card">
			<Table.Row>
				<Table.Head>Action</Table.Head>
				<Table.Head class="w-px">Operation</Table.Head>
				<Table.Head>Path</Table.Head>
				<Table.Head>Worker</Table.Head>
				<Table.Head class="w-px text-right">Run</Table.Head>
			</Table.Row>
		</Table.Header>
		<Table.Body>
			{#if actions.length > 0}
				{#each actions as action (action.id)}
					{@const reason = disabledReason(action)}
					<Table.Row>
						<Table.Cell>
							<div class="grid gap-1">
								<strong>{action.label}</strong>
								<span class="text-xs text-muted-foreground">{action.description}</span>
							</div>
						</Table.Cell>
						<Table.Cell class="w-px">
							<Badge variant="outline">{action.operation}</Badge>
						</Table.Cell>
						<Table.Cell class="whitespace-nowrap font-mono text-xs">
							{action.method}
							{action.path}
						</Table.Cell>
						<Table.Cell class="whitespace-nowrap text-xs">{action.workerPath}</Table.Cell>
						<Table.Cell class="w-px text-right">
							<Tooltip.Root>
								<Tooltip.Trigger>
									{#snippet child({ props })}
										<Button
											{...props}
											type="button"
											size="icon-sm"
											variant="outline"
											disabled={!!reason}
											aria-label={`Run ${action.label}`}
										>
											<CirclePlayIcon aria-hidden="true" />
										</Button>
									{/snippet}
								</Tooltip.Trigger>
								<Tooltip.Content>{reason || `Run ${action.label}`}</Tooltip.Content>
							</Tooltip.Root>
						</Table.Cell>
					</Table.Row>
				{/each}
			{:else if loading}
				<Table.Row>
					<Table.Cell colspan={5} class="text-muted-foreground">Loading manual actions.</Table.Cell>
				</Table.Row>
			{:else if errorMessage}
				<Table.Row>
					<Table.Cell colspan={5} class="text-destructive">{errorMessage}</Table.Cell>
				</Table.Row>
			{:else}
				<Table.Row>
					<Table.Cell colspan={5} class="text-muted-foreground"
						>No manual actions registered.</Table.Cell
					>
				</Table.Row>
			{/if}
		</Table.Body>
	</Table.Root>
</div>
