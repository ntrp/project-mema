<script lang="ts">
	import SettingsFormModal from '$lib/components/settings/shared/SettingsFormModal.svelte';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import { Textarea } from '$lib/components/ui/textarea';
	import type { DLNARendererProfile } from '$lib/settings/types';

	interface Props {
		cloneSource?: DLNARendererProfile;
		cloneId: string;
		cloneName: string;
		importOpen: boolean;
		importText: string;
		saving?: boolean;
		onCloneId: (_value: string) => void;
		onCloneName: (_value: string) => void;
		onImportText: (_value: string) => void;
		onCloseClone: () => void;
		onCloseImport: () => void;
		onClone: () => void | Promise<void>;
		onImport: () => void | Promise<void>;
	}

	let {
		cloneSource,
		cloneId,
		cloneName,
		importOpen,
		importText,
		saving = false,
		onCloneId,
		onCloneName,
		onImportText,
		onCloseClone,
		onCloseImport,
		onClone,
		onImport
	}: Props = $props();
</script>

{#if cloneSource}
	<SettingsFormModal title={`Clone ${cloneSource.name}`} onClose={onCloseClone}>
		<form
			class="grid gap-4"
			onsubmit={(event) => {
				event.preventDefault();
				void onClone();
			}}
		>
			<div class="grid gap-2">
				<Label for="dlna-clone-id">New profile ID</Label>
				<Input
					id="dlna-clone-id"
					value={cloneId}
					oninput={(event) => onCloneId(event.currentTarget.value)}
					required
				/>
			</div>
			<div class="grid gap-2">
				<Label for="dlna-clone-name">New profile name</Label>
				<Input
					id="dlna-clone-name"
					value={cloneName}
					oninput={(event) => onCloneName(event.currentTarget.value)}
					required
				/>
			</div>
			<div class="flex justify-end gap-2">
				<Button type="button" variant="outline" onclick={onCloseClone}>Cancel</Button>
				<Button type="submit" disabled={saving}>{saving ? 'Cloning' : 'Clone profile'}</Button>
			</div>
		</form>
	</SettingsFormModal>
{/if}

{#if importOpen}
	<SettingsFormModal title="Import DLNA profile" onClose={onCloseImport}>
		<form
			class="grid gap-4"
			onsubmit={(event) => {
				event.preventDefault();
				void onImport();
			}}
		>
			<div class="grid gap-2">
				<Label for="dlna-profile-import">Profile JSON</Label>
				<Textarea
					id="dlna-profile-import"
					class="min-h-80 font-mono text-xs"
					value={importText}
					oninput={(event) => onImportText(event.currentTarget.value)}
					required
					spellcheck={false}
				/>
			</div>
			<div class="flex justify-end gap-2">
				<Button type="button" variant="outline" onclick={onCloseImport}>Cancel</Button>
				<Button type="submit" disabled={saving}>{saving ? 'Importing' : 'Import profile'}</Button>
			</div>
		</form>
	</SettingsFormModal>
{/if}
