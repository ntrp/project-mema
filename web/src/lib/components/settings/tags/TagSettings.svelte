<script lang="ts">
	import PlusIcon from '@lucide/svelte/icons/plus';
	import SettingsFormModal from '$lib/components/settings/shared/SettingsFormModal.svelte';
	import SettingsRowActionButton from '$lib/components/settings/shared/SettingsRowActionButton.svelte';
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import { Card } from '$lib/components/ui/card';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import * as Table from '$lib/components/ui/table';
	import { formatDate } from '$lib/settings/dateFormat';
	import type { Tag, TagForm } from '$lib/settings/types';

	interface Props {
		tags: Tag[];
		form: TagForm;
		saving: boolean;
		deletingId?: string;
		onSave: (_event: SubmitEvent) => void | Promise<void>;
		onCancel: () => void;
		onEdit: (_tag: Tag) => void;
		onDelete: (_id: string) => void | Promise<void>;
	}

	let {
		tags,
		form = $bindable(),
		saving,
		deletingId,
		onSave,
		onCancel,
		onEdit,
		onDelete
	}: Props = $props();

	let tagModalOpen = $state(false);

	function openTagModal() {
		onCancel();
		tagModalOpen = true;
	}

	function editTag(tag: Tag) {
		onEdit(tag);
		tagModalOpen = true;
	}

	function closeTagModal() {
		onCancel();
		tagModalOpen = false;
	}

	async function saveTag(event: SubmitEvent) {
		await onSave(event);
		if (!form.id && form.name === '') {
			tagModalOpen = false;
		}
	}
</script>

<Card class="gap-0 p-0" aria-label="Tags">
	<div class="flex justify-end border-b px-4 py-3">
		<Button type="button" onclick={openTagModal}>
			<PlusIcon aria-hidden="true" />
			<span>Add tag</span>
		</Button>
	</div>

	<Table.Root>
		<Table.Header>
			<Table.Row>
				<Table.Head>Name</Table.Head>
				<Table.Head>Updated</Table.Head>
				<Table.Head class="text-right">Actions</Table.Head>
			</Table.Row>
		</Table.Header>
		<Table.Body>
			{#each tags as tag (tag.id)}
				<Table.Row>
					<Table.Cell><Badge variant="secondary">{tag.name}</Badge></Table.Cell>
					<Table.Cell>{formatDate(tag.updatedAt)}</Table.Cell>
					<Table.Cell>
						<div class="flex justify-end gap-2">
							<SettingsRowActionButton
								label={`Edit ${tag.name}`}
								icon="edit"
								onclick={() => editTag(tag)}
							/>
							<SettingsRowActionButton
								label={`${deletingId === tag.id ? 'Deleting' : 'Delete'} ${tag.name}`}
								icon="delete"
								variant="destructive"
								disabled={deletingId === tag.id}
								confirmTitle="Delete tag"
								confirmDescription={`Delete tag "${tag.name}"?`}
								confirmLabel="Delete tag"
								onclick={() => onDelete(tag.id)}
							/>
						</div>
					</Table.Cell>
				</Table.Row>
			{:else}
				<Table.Row>
					<Table.Cell colspan={3} class="py-8 text-center text-muted-foreground">
						No tags configured
					</Table.Cell>
				</Table.Row>
			{/each}
		</Table.Body>
	</Table.Root>

	{#if tagModalOpen}
		<SettingsFormModal title={form.id ? 'Edit tag' : 'Add tag'} onClose={closeTagModal}>
			<form class="grid gap-4" onsubmit={saveTag}>
				<div class="grid gap-2">
					<Label for="tag-name">Name</Label>
					<Input id="tag-name" bind:value={form.name} type="text" maxlength={80} required />
				</div>
				<div class="flex justify-end gap-2">
					<Button type="button" variant="outline" onclick={closeTagModal}>Cancel</Button>
					<Button type="submit" disabled={saving}>
						{saving ? 'Saving' : form.id ? 'Update tag' : 'Create tag'}
					</Button>
				</div>
			</form>
		</SettingsFormModal>
	{/if}
</Card>
