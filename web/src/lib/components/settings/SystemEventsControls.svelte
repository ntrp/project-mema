<script lang="ts">
	const severityOptions = ['info', 'warning', 'error'] as const;
	type SeverityFilter = (typeof severityOptions)[number];

	interface Props {
		severityFilter: SeverityFilter;
		loading: boolean;
		clearing: boolean;
		eventsEmpty: boolean;
		onSeverityChange: (severity: SeverityFilter) => void;
		onClear: () => void;
	}

	let { severityFilter, loading, clearing, eventsEmpty, onSeverityChange, onClear }: Props =
		$props();

	function severityLabel(severity: SeverityFilter) {
		return severity[0].toUpperCase() + severity.slice(1);
	}

	function handleSeverityChange(event: Event) {
		const target = event.currentTarget as unknown as { value: string };
		onSeverityChange(target.value as SeverityFilter);
	}
</script>

<div class="log-controls events-controls">
	<label class="severity-filter">
		<span>Severity</span>
		<select value={severityFilter} onchange={handleSeverityChange}>
			{#each severityOptions as severity (severity)}
				<option value={severity}>{severityLabel(severity)}</option>
			{/each}
		</select>
	</label>
	<button
		type="button"
		class="danger icon-button events-clear-button"
		aria-label={clearing ? 'Clearing events' : 'Clear all events'}
		title={clearing ? 'Clearing events' : 'Clear all events'}
		disabled={loading || clearing || eventsEmpty}
		onclick={onClear}
	>
		<span class="app-icon" aria-hidden="true">delete</span>
	</button>
</div>
