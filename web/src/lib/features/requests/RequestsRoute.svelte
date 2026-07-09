<script lang="ts">
	import MediaRequestArea from '$lib/components/app/requests/MediaRequestArea.svelte';
	import { getAppShellContext } from '$lib/features/app/appShellContext';

	interface Props {
		selectedRequestId?: string;
	}

	let { selectedRequestId }: Props = $props();
	const app = getAppShellContext();
</script>

<MediaRequestArea
	requests={app.mediaRequests}
	{selectedRequestId}
	libraryFolders={app.libraryFolders}
	qualityProfiles={app.mediaProfiles}
	canManage={app.isAdmin}
	approvingRequestId={app.approvingRequestId}
	onApprove={app.approveMediaRequest}
/>
