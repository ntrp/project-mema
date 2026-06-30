<script lang="ts">
	import { cancellable, deletable, manualImportable } from './activityDisplay';
	import type { DownloadActivity } from '$lib/settings/types';

	interface Props {
		activity: DownloadActivity;
		canManage: boolean;
		cancellingId?: string;
		deletingId?: string;
		onManualImport: (_activity: DownloadActivity) => void;
		onCancel: (_activity: DownloadActivity) => void;
		onDelete: (_activity: DownloadActivity) => void;
	}

	let { activity, canManage, cancellingId, deletingId, onManualImport, onCancel, onDelete }: Props =
		$props();
</script>

{#if canManage && manualImportable(activity)}
	<button
		type="button"
		class="secondary icon-button"
		aria-label={`Manual import ${activity.releaseTitle}`}
		title="Manual import"
		onclick={() => onManualImport(activity)}
	>
		<span class="app-icon" aria-hidden="true">upload_file</span>
	</button>
{/if}
{#if canManage && cancellable(activity)}
	<button
		type="button"
		class="danger icon-button"
		aria-label={`Cancel ${activity.releaseTitle}`}
		title="Cancel"
		disabled={cancellingId === activity.id}
		onclick={() => onCancel(activity)}
	>
		<span class="app-icon" aria-hidden="true">
			{cancellingId === activity.id ? 'sync' : 'close'}
		</span>
	</button>
{/if}
{#if canManage && deletable(activity)}
	<button
		type="button"
		class="danger icon-button"
		aria-label={`Delete ${activity.releaseTitle}`}
		title="Delete"
		disabled={deletingId === activity.id}
		onclick={() => onDelete(activity)}
	>
		<span class="app-icon" aria-hidden="true">
			{deletingId === activity.id ? 'sync' : 'delete'}
		</span>
	</button>
{/if}
