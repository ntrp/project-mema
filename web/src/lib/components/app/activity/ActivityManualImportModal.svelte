<script lang="ts">
	import SettingsFormModal from '$lib/components/settings/shared/SettingsFormModal.svelte';
	import ActivityManualImportForm from './ActivityManualImportForm.svelte';
	import type { DownloadActivity, ManualImportRequest } from '$lib/settings/types';

	interface Props {
		activity: DownloadActivity;
		importing: boolean;
		error?: string;
		onImport: (_request: ManualImportRequest) => void;
		onClose: () => void;
	}

	let { activity, importing, error, onImport, onClose }: Props = $props();
</script>

<SettingsFormModal
	title={activity.mediaTitle}
	modalClass="w-[min(1840px,calc(100vw-48px))]"
	{onClose}
>
	<ActivityManualImportForm {activity} {importing} {error} {onImport} {onClose} />
</SettingsFormModal>
