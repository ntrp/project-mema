<script lang="ts">
	import { Checkbox } from '$lib/components/ui/checkbox';
	import * as Table from '$lib/components/ui/table';

	interface Props {
		checked: boolean;
		disabled: boolean;
		showSeriesControls?: boolean;
		onToggle: () => void;
	}

	let { checked, disabled, showSeriesControls = true, onToggle }: Props = $props();
</script>

<Table.Header class="sticky top-0 bg-card">
	<Table.Row>
		<Table.Head class="w-px">
			<Checkbox
				{checked}
				{disabled}
				aria-label="Select visible importable rows"
				onclick={onToggle}
			/>
		</Table.Head>
		<Table.Head>Directory / File</Table.Head>
		<Table.Head class="w-px whitespace-nowrap">Metadata provider</Table.Head>
		<Table.Head class="w-px whitespace-nowrap">Quality profile</Table.Head>
		<Table.Head class="w-px whitespace-nowrap">Monitor</Table.Head>
		<Table.Head class="w-px whitespace-nowrap">Minimum Availability</Table.Head>
		<Table.Head class="w-px whitespace-nowrap">{showSeriesControls ? 'Series type' : ''}</Table.Head>
	</Table.Row>
</Table.Header>
