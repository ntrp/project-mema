<script lang="ts">
	import PlusIcon from '@lucide/svelte/icons/plus';
	import TrashIcon from '@lucide/svelte/icons/trash-2';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import * as Select from '$lib/components/ui/select';
	import * as Table from '$lib/components/ui/table';
	import type { SubtitleProviderForm } from '$lib/settings/types';

	type MockRows = NonNullable<SubtitleProviderForm['mockSubtitles']>;
	type MockRow = MockRows[number];

	interface Props {
		rows: MockRows;
		onChange: (_rows: MockRows) => void;
	}

	let { rows, onChange }: Props = $props();

	function addRow() {
		onChange([...rows, { title: '', languageId: 'english', format: 'srt' }]);
	}

	function updateRow(index: number, value: Partial<MockRow>) {
		onChange(rows.map((row, rowIndex) => (rowIndex === index ? { ...row, ...value } : row)));
	}

	function removeRow(index: number) {
		onChange(rows.filter((_, rowIndex) => rowIndex !== index));
	}
</script>

<div class="grid gap-2">
	<div class="flex items-center justify-between gap-3">
		<span class="text-sm font-bold text-muted-foreground">Mock subtitles</span>
		<Button type="button" variant="outline" size="sm" onclick={addRow}>
			<PlusIcon class="size-4" />
			Add row
		</Button>
	</div>
	<Table.Root class="w-full table-fixed">
		<Table.Header>
			<Table.Row>
				<Table.Head class="text-left">Title</Table.Head>
				<Table.Head class="w-40 text-left">Language</Table.Head>
				<Table.Head class="w-32 text-left">Format</Table.Head>
				<Table.Head class="w-12 text-right"><span class="sr-only">Actions</span></Table.Head>
			</Table.Row>
		</Table.Header>
		<Table.Body>
			{#each rows as row, index (index)}
				<Table.Row>
					<Table.Cell>
						<Input
							value={row.title}
							required
							maxlength={500}
							aria-label="Mock subtitle title"
							oninput={(event) => updateRow(index, { title: event.currentTarget.value })}
						/>
					</Table.Cell>
					<Table.Cell>
						<Input
							value={row.languageId}
							required
							maxlength={100}
							aria-label="Mock subtitle language"
							oninput={(event) => updateRow(index, { languageId: event.currentTarget.value })}
						/>
					</Table.Cell>
					<Table.Cell>
						<Select.Root
							type="single"
							value={row.format}
							onValueChange={(format) => updateRow(index, { format })}
						>
							<Select.Trigger class="w-full">{row.format.toUpperCase()}</Select.Trigger>
							<Select.Content>
								<Select.Item value="srt" label="SRT" />
								<Select.Item value="vtt" label="VTT" />
								<Select.Item value="ass" label="ASS" />
								<Select.Item value="ssa" label="SSA" />
							</Select.Content>
						</Select.Root>
					</Table.Cell>
					<Table.Cell class="text-right">
						<Button
							type="button"
							variant="ghost"
							size="icon"
							aria-label="Remove mock subtitle"
							onclick={() => removeRow(index)}
						>
							<TrashIcon class="size-4" />
						</Button>
					</Table.Cell>
				</Table.Row>
			{:else}
				<Table.Row>
					<Table.Cell colspan={4} class="py-6 text-center text-muted-foreground">
						No mock subtitles configured.
					</Table.Cell>
				</Table.Row>
			{/each}
		</Table.Body>
	</Table.Root>
</div>
