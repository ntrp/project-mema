<script lang="ts">
	import ConfirmActionButton from '$lib/components/shared/ConfirmActionButton.svelte';
	import IntegrationTestStatus from '$lib/components/settings/shared/IntegrationTestStatus.svelte';
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import type {
		IntegrationTestResults,
		SubtitleProvider,
		SubtitleProviderCatalogEntry
	} from '$lib/settings/types';
	import { runtimeLabel } from '../catalog/subtitleProviderCatalogFilters';

	interface Props {
		providers: SubtitleProvider[];
		catalog: SubtitleProviderCatalogEntry[];
		testingId?: string;
		testResults: IntegrationTestResults;
		onEdit: (_provider: SubtitleProvider) => void;
		onDelete: (_id: string) => void | Promise<void>;
		onTest: (_id: string) => void | Promise<void>;
	}

	let { providers, catalog, testingId, testResults, onEdit, onDelete, onTest }: Props = $props();

	function entryFor(provider: SubtitleProvider) {
		return catalog.find(
			(entry) => entry.key === provider.catalogKey || entry.key === provider.type
		);
	}
</script>

<div class="overflow-hidden rounded-md border border-border">
	<table class="w-full min-w-180 table-auto border-collapse text-sm">
		<thead class="bg-card text-left text-xs font-extrabold text-muted-foreground">
			<tr class="border-b border-border">
				<th class="px-3 py-2">Name</th>
				<th class="px-3 py-2">Support</th>
				<th class="px-3 py-2">Enabled</th>
				<th class="px-3 py-2">Test</th>
				<th class="w-px px-3 py-2">Actions</th>
			</tr>
		</thead>
		<tbody>
			{#each providers as provider (provider.id)}
				{@const entry = entryFor(provider)}
				{@const supported = provider.runtimeStatus === 'supported'}
				<tr class="border-b border-border last:border-0">
					<td class="px-3 py-2">
						<div class="font-bold text-foreground">{provider.name}</div>
						<div class="text-xs text-muted-foreground">{entry?.displayName ?? provider.type}</div>
					</td>
					<td class="px-3 py-2">
						<Badge variant="outline">{runtimeLabel(provider.runtimeStatus)}</Badge>
					</td>
					<td class="px-3 py-2">{provider.enabled ? 'Enabled' : 'Disabled'}</td>
					<td class="px-3 py-2">
						<IntegrationTestStatus
							enabled={provider.enabled}
							result={testResults[provider.id]}
							testing={testingId === provider.id}
						/>
					</td>
					<td class="px-3 py-2">
						<div class="flex justify-end gap-2">
							<Button type="button" size="sm" variant="outline" onclick={() => onEdit(provider)}>
								Edit
							</Button>
							<Button
								type="button"
								size="sm"
								variant="outline"
								disabled={!supported || testingId === provider.id}
								onclick={() => onTest(provider.id)}
							>
								Test
							</Button>
							<ConfirmActionButton
								label={`Delete ${provider.name}`}
								title="Delete subtitle provider"
								description={`Delete subtitle provider "${provider.name}"?`}
								confirmLabel="Delete provider"
								onConfirm={() => onDelete(provider.id)}
							>
								Delete
							</ConfirmActionButton>
						</div>
					</td>
				</tr>
			{:else}
				<tr>
					<td class="px-3 py-8 text-center text-muted-foreground" colspan="5">
						No subtitle providers configured.
					</td>
				</tr>
			{/each}
		</tbody>
	</table>
</div>
