<script lang="ts">
	import SettingsFormModal from '$lib/components/settings/shared/SettingsFormModal.svelte';
	import { Button } from '$lib/components/ui/button';
	import type { Language } from '$lib/settings/types';

	interface Props {
		language: Language;
		deleting: boolean;
		onCancel: () => void;
		onConfirm: () => void | Promise<void>;
	}

	let { language, deleting, onCancel, onConfirm }: Props = $props();
</script>

<SettingsFormModal
	title="Delete language"
	modalClass="grid w-[min(460px,100%)] gap-4"
	onClose={onCancel}
>
	<div class="grid gap-[18px]">
		<p class="m-0 leading-normal text-muted-foreground">
			Delete {language.displayName} ({language.code}) from the release language catalog.
		</p>
		<div class="flex flex-wrap justify-end gap-2.5">
			<Button type="button" variant="outline" disabled={deleting} onclick={onCancel}>Cancel</Button>
			<Button type="button" variant="destructive" disabled={deleting} onclick={onConfirm}>
				{deleting ? 'Deleting' : 'Delete language'}
			</Button>
		</div>
	</div>
</SettingsFormModal>
