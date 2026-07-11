<script lang="ts">
	import PersonDetailArea from '$lib/components/app/media/person-detail/PersonDetailArea.svelte';
	import { getAppShellContext } from '$lib/features/app/appShellContext';
	import { createMediaItemsQuery } from '$lib/features/library/queries.svelte';

	const app = getAppShellContext();
	const library = createMediaItemsQuery();
</script>

<PersonDetailArea
	person={app.personDetail}
	loading={app.loadingPersonDetail}
	mediaItems={library.data ?? []}
	addingKey={app.addingKey}
	actionLabel={app.isAdmin ? 'Add' : 'Request'}
	onAdd={app.addMedia}
/>
