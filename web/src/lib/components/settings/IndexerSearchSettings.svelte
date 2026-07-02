<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import type { IndexerSearchResponse, IndexerSearchSettings } from '$lib/settings/types';

	interface Props {
		search: IndexerSearchResponse;
		clearing: boolean;
		saving: boolean;
		onClearCache: () => void | Promise<void>;
		onSaveSettings: (_settings: IndexerSearchSettings) => void | Promise<void>;
	}

	let { search, clearing, saving, onClearCache, onSaveSettings }: Props = $props();
	let cacheDurationMinutes = $state(0);
	let historyRetentionDays = $state(30);
	const settingsChanged = $derived(
		cacheDurationMinutes !== search.settings.cacheDurationMinutes ||
			historyRetentionDays !== search.settings.historyRetentionDays
	);

	$effect(() => {
		cacheDurationMinutes = search.settings.cacheDurationMinutes;
		historyRetentionDays = search.settings.historyRetentionDays;
	});

	function saveSettings() {
		void onSaveSettings({ cacheDurationMinutes, historyRetentionDays });
	}
</script>

<Card.Root aria-labelledby="indexer-search-settings-title">
	<Card.Header>
		<div>
			<Card.Description>Search</Card.Description>
			<Card.Title id="indexer-search-settings-title">Indexer search settings</Card.Title>
		</div>
	</Card.Header>
	<Card.Content class="grid gap-4">
		<div class="grid items-end gap-3 md:grid-cols-[minmax(0,1fr)_minmax(0,1fr)_auto]">
			<div class="grid gap-1.5">
				<Label for="indexer-cache-duration">Cache duration minutes</Label>
				<Input
					id="indexer-cache-duration"
					type="number"
					min="0"
					max="43200"
					bind:value={cacheDurationMinutes}
				/>
			</div>
			<div class="grid gap-1.5">
				<Label for="indexer-history-retention">History retention days</Label>
				<Input
					id="indexer-history-retention"
					type="number"
					min="1"
					max="365"
					bind:value={historyRetentionDays}
				/>
			</div>
			<div class="flex flex-wrap justify-end gap-2">
				<Button type="button" disabled={saving || !settingsChanged} onclick={saveSettings}>
					{saving ? 'Saving' : 'Save settings'}
				</Button>
				<Button
					type="button"
					variant="destructive"
					disabled={clearing}
					onclick={() => void onClearCache()}
				>
					{clearing ? 'Resetting' : 'Reset cache'}
				</Button>
			</div>
		</div>
	</Card.Content>
</Card.Root>
