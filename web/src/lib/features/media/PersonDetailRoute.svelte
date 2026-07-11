<script lang="ts">
	import PersonDetailArea from '$lib/components/app/media/person-detail/PersonDetailArea.svelte';
	import { getAppShellContext } from '$lib/features/app/appShellContext';
	import { createMediaItemsQuery } from '$lib/features/library/queries.svelte';
	import { createPersonDetailQuery } from './queries.svelte';

	const app = getAppShellContext();
	const library = createMediaItemsQuery();
	const person = createPersonDetailQuery(
		() => app.route?.personProvider,
		() => app.route?.personId
	);
</script>

<PersonDetailArea
	person={person.data}
	loading={person.isFetching}
	mediaItems={library.data ?? []}
	addingKey={app.addingKey}
	actionLabel={app.isAdmin ? 'Add' : 'Request'}
	onAdd={app.addMedia}
/>
