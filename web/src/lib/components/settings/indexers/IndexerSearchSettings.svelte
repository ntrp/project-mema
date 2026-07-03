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
	let cacheDurationMinutesDraft = $state<number | undefined>();
	let historyRetentionDaysDraft = $state<number | undefined>();
	let automaticBlocklistExpiryDaysDraft = $state<number | undefined>();
	const cacheDurationMinutes = $derived(
		cacheDurationMinutesDraft ?? search.settings.cacheDurationMinutes
	);
	const historyRetentionDays = $derived(
		historyRetentionDaysDraft ?? search.settings.historyRetentionDays
	);
	const automaticBlocklistExpiryDays = $derived(
		automaticBlocklistExpiryDaysDraft ?? search.settings.automaticBlocklistExpiryDays
	);
	const settingsChanged = $derived(
		cacheDurationMinutes !== search.settings.cacheDurationMinutes ||
			historyRetentionDays !== search.settings.historyRetentionDays ||
			automaticBlocklistExpiryDays !== search.settings.automaticBlocklistExpiryDays
	);

	function saveSettings() {
		void Promise.resolve(
			onSaveSettings({
				cacheDurationMinutes,
				historyRetentionDays,
				automaticBlocklistExpiryDays
			})
		).then(() => {
			cacheDurationMinutesDraft = undefined;
			historyRetentionDaysDraft = undefined;
			automaticBlocklistExpiryDaysDraft = undefined;
		});
	}

	function numberFromInput(event: Event) {
		return Number((event.currentTarget as HTMLInputElement).value);
	}

	function updateCacheDuration(event: Event) {
		cacheDurationMinutesDraft = numberFromInput(event);
	}

	function updateHistoryRetention(event: Event) {
		historyRetentionDaysDraft = numberFromInput(event);
	}

	function updateAutomaticBlocklistExpiry(event: Event) {
		automaticBlocklistExpiryDaysDraft = numberFromInput(event);
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
		<div class="grid items-end gap-3 md:grid-cols-[repeat(3,minmax(0,1fr))_auto]">
			<div class="grid gap-1.5">
				<Label for="indexer-cache-duration">Cache duration minutes</Label>
				<Input
					id="indexer-cache-duration"
					type="number"
					min="0"
					max="43200"
					value={cacheDurationMinutes}
					oninput={updateCacheDuration}
				/>
			</div>
			<div class="grid gap-1.5">
				<Label for="indexer-history-retention">History retention days</Label>
				<Input
					id="indexer-history-retention"
					type="number"
					min="1"
					max="365"
					value={historyRetentionDays}
					oninput={updateHistoryRetention}
				/>
			</div>
			<div class="grid gap-1.5">
				<Label for="automatic-blocklist-expiry">Automatic blocklist expiry days</Label>
				<Input
					id="automatic-blocklist-expiry"
					type="number"
					min="1"
					max="365"
					value={automaticBlocklistExpiryDays}
					oninput={updateAutomaticBlocklistExpiry}
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
