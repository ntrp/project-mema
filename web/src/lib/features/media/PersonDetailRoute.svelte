<script lang="ts">
	import PersonDetailArea from '$lib/components/app/media/person-detail/PersonDetailArea.svelte';
	import { getAppShellContext } from '$lib/features/app/appShellContext';

	const app = getAppShellContext();
</script>

<PersonDetailArea
	person={app.personDetail}
	loading={app.loadingPersonDetail}
	mediaItems={app.mediaItems}
	addingKey={app.addingKey}
	actionLabel={app.isAdmin ? 'Add' : 'Request'}
	onAdd={app.addMedia}
/>
