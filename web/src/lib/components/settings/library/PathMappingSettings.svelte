<script lang="ts">
	import PlusIcon from '@lucide/svelte/icons/plus';
	import SettingsRowActionButton from '$lib/components/settings/shared/SettingsRowActionButton.svelte';
	import SettingsFormModal from '$lib/components/settings/shared/SettingsFormModal.svelte';
	import { Button } from '$lib/components/ui/button';
	import { Card } from '$lib/components/ui/card';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import { Separator } from '$lib/components/ui/separator';
	import * as Table from '$lib/components/ui/table';
	import type { PathMapping, PathMappingForm } from '$lib/settings/types';

	interface Props {
		mappings: PathMapping[];
		form: PathMappingForm;
		saving: boolean;
		deletingId?: string;
		onSave: (_event: SubmitEvent) => void | Promise<void>;
		onDelete: (_id: string) => void | Promise<void>;
	}

	let { mappings, form = $bindable(), saving, deletingId, onSave, onDelete }: Props = $props();
	let modalOpen = $state(false);

	async function save(event: SubmitEvent) {
		await onSave(event);
		if (form.clientPath.trim() === '' && form.appPath.trim() === '') {
			modalOpen = false;
		}
	}
</script>

<section class="grid gap-4" aria-labelledby="path-mapping-title">
	<div class="grid gap-2">
		<div class="flex items-center justify-between gap-3">
			<h3 id="path-mapping-title" class="m-0 text-lg text-foreground">Path Mappings</h3>
			<Button type="button" onclick={() => (modalOpen = true)}>
				<PlusIcon aria-hidden="true" />
				<span>Add path</span>
			</Button>
		</div>
		<Separator />
	</div>
	<Card class="p-0">
		<Table.Root>
			<Table.Header>
				<Table.Row>
					<Table.Head>Client path</Table.Head>
					<Table.Head>App path</Table.Head>
					<Table.Head class="text-right">Actions</Table.Head>
				</Table.Row>
			</Table.Header>
			<Table.Body>
				{#each mappings as mapping (mapping.id)}
					<Table.Row>
						<Table.Cell class="max-w-[360px] truncate">{mapping.clientPath}</Table.Cell>
						<Table.Cell class="max-w-[360px] truncate">{mapping.appPath}</Table.Cell>
						<Table.Cell>
							<div class="flex justify-end">
								<SettingsRowActionButton
									label="Delete path mapping"
									icon="delete"
									variant="destructive"
									disabled={deletingId === mapping.id}
									confirmTitle="Delete path mapping"
									confirmDescription={`Delete path mapping from "${mapping.clientPath}" to "${mapping.appPath}"?`}
									confirmLabel="Delete mapping"
									onclick={() => onDelete(mapping.id)}
								/>
							</div>
						</Table.Cell>
					</Table.Row>
				{:else}
					<Table.Row>
						<Table.Cell colspan={3} class="py-8 text-center text-muted-foreground">
							No paths have been defined.
						</Table.Cell>
					</Table.Row>
				{/each}
			</Table.Body>
		</Table.Root>
	</Card>
	{#if modalOpen}
		<SettingsFormModal title="Add path mapping" onClose={() => (modalOpen = false)}>
			<form class="grid gap-4 sm:grid-cols-2" onsubmit={save}>
				<div class="grid gap-2">
					<Label for="path-mapping-client-path">Client path</Label>
					<Input
						id="path-mapping-client-path"
						bind:value={form.clientPath}
						placeholder="/downloads"
						required
					/>
				</div>
				<div class="grid gap-2">
					<Label for="path-mapping-app-path">App path</Label>
					<Input
						id="path-mapping-app-path"
						bind:value={form.appPath}
						placeholder="/mnt/downloads"
						required
					/>
				</div>
				<div class="flex justify-end sm:col-span-2">
					<Button type="submit" disabled={saving}>
						{saving ? 'Saving' : 'Save mapping'}
					</Button>
				</div>
			</form>
		</SettingsFormModal>
	{/if}
</section>
