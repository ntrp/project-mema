<script lang="ts">
	import TrashIcon from '@lucide/svelte/icons/trash-2';
	import SettingsSelect from '$lib/components/settings/shared/SettingsSelect.svelte';
	import { Button } from '$lib/components/ui/button';
	import * as Tooltip from '$lib/components/ui/tooltip';
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

	const severitySelectOptions = severityOptions.map((severity) => ({
		value: severity,
		label: severityLabel(severity)
	}));
</script>

<div class="flex flex-wrap items-center justify-end gap-3">
	<SettingsSelect
		value={severityFilter}
		options={severitySelectOptions}
		onValueChange={(value) => onSeverityChange(value as SeverityFilter)}
	/>
	<Tooltip.Root>
		<Tooltip.Trigger>
			{#snippet child({ props })}
				<Button
					{...props}
					type="button"
					variant="destructive"
					size="icon-sm"
					aria-label={clearing ? 'Clearing events' : 'Clear all events'}
					disabled={loading || clearing || eventsEmpty}
					onclick={onClear}
				>
					<TrashIcon aria-hidden="true" />
				</Button>
			{/snippet}
		</Tooltip.Trigger>
		<Tooltip.Content>{clearing ? 'Clearing events' : 'Clear all events'}</Tooltip.Content>
	</Tooltip.Root>
</div>
