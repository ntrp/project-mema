<script lang="ts">
	import { onMount } from 'svelte';
	import NoticeStack from '$lib/components/settings/shared/NoticeStack.svelte';
	import SectionHeading from '$lib/components/shared/SectionHeading.svelte';
	import { Button } from '$lib/components/ui/button';
	import { Card } from '$lib/components/ui/card';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import * as Select from '$lib/components/ui/select';
	import { createFileDeleteResource } from '$lib/features/settings/resources/filePolicies.svelte';
	import type { FileDeleteMode, FileDeleteSettingsRequest } from '$lib/settings/types';

	const modeOptions: { value: FileDeleteMode; label: string; description: string }[] = [
		{ value: 'permanent', label: 'Delete permanently', description: 'Remove files from disk.' },
		{
			value: 'recycle',
			label: 'Move to recycle folder',
			description: 'Move files under each root.'
		},
		{ value: 'keep', label: 'Keep files', description: 'Record a skipped delete event.' }
	];

	let form = $state<FileDeleteSettingsRequest>({ mode: 'permanent', recycleFolder: '.recycle' });
	const resource = createFileDeleteResource();
	let message = $state('');
	let errorMessage = $state('');
	const selectedMode = $derived(modeOptions.find((option) => option.value === form.mode));
	const recycleInvalid = $derived(
		form.recycleFolder.trim() === '' ||
			form.recycleFolder.startsWith('/') ||
			form.recycleFolder.includes('..') ||
			!form.recycleFolder.split('/')[0]?.startsWith('.')
	);

	const loading = $derived(resource.query.isPending || resource.query.isFetching);
	const saving = $derived(resource.save.isPending);
	onMount(() => {
		void loadSettings();
	});

	async function loadSettings() {
		message = '';
		errorMessage = '';
		try {
			const { data: settings } = await resource.query.refetch();
			if (!settings) return;
			form = { mode: settings.mode, recycleFolder: settings.recycleFolder };
		} catch (error) {
			errorMessage = error instanceof Error ? error.message : 'Could not load file delete settings';
		} finally {
			/* query owns loading state */
		}
	}

	async function saveSettings(event: SubmitEvent) {
		event.preventDefault();
		message = '';
		errorMessage = '';
		if (recycleInvalid) {
			errorMessage = 'Recycle folder must be a hidden relative folder';
			return;
		}
		try {
			const settings = await resource.save.mutateAsync({
				mode: form.mode,
				recycleFolder: form.recycleFolder.trim()
			});
			form = { mode: settings.mode, recycleFolder: settings.recycleFolder };
			message = 'File delete settings saved';
		} catch (error) {
			errorMessage = error instanceof Error ? error.message : 'Could not save file delete settings';
		} finally {
			/* mutation owns saving state */
		}
	}
</script>

<Card class="overflow-visible p-5" aria-label="File delete policy">
	<form class="grid gap-4" onsubmit={saveSettings}>
		<SectionHeading title="File Delete Policy">
			{#snippet actions()}
				<div class="flex flex-wrap justify-end gap-2.5">
					<Button
						type="button"
						variant="outline"
						disabled={loading || saving}
						onclick={loadSettings}
					>
						Reload
					</Button>
					<Button type="submit" disabled={loading || saving || recycleInvalid}>
						{saving ? 'Saving' : 'Save policy'}
					</Button>
				</div>
			{/snippet}
		</SectionHeading>

		<NoticeStack {message} {errorMessage} />

		<div class="grid gap-4 sm:grid-cols-[minmax(0,1fr)_minmax(0,1fr)]">
			<div class="space-y-2">
				<Label for="file-delete-mode">Mode</Label>
				<Select.Root
					type="single"
					value={form.mode}
					onValueChange={(value) => (form = { ...form, mode: value as FileDeleteMode })}
				>
					<Select.Trigger id="file-delete-mode" class="w-full">
						{selectedMode?.label ?? 'Select mode'}
					</Select.Trigger>
					<Select.Content>
						{#each modeOptions as option (option.value)}
							<Select.Item value={option.value} label={option.label} />
						{/each}
					</Select.Content>
				</Select.Root>
				<p class="m-0 text-sm text-muted-foreground">{selectedMode?.description}</p>
			</div>
			<div class="space-y-2">
				<Label for="file-recycle-folder">Recycle folder</Label>
				<Input
					id="file-recycle-folder"
					bind:value={form.recycleFolder}
					aria-invalid={recycleInvalid}
				/>
				<p class="m-0 text-sm text-muted-foreground">
					Relative hidden folder under each library root.
				</p>
			</div>
		</div>
	</form>
</Card>
