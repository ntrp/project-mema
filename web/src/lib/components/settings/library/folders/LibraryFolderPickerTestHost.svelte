<script lang="ts">
	import { QueryClientProvider } from '@tanstack/svelte-query';
	import { createAppQueryClient } from '$lib/app/query/queryClient';
	import LibraryFolderPicker from './LibraryFolderPicker.svelte';

	interface Props {
		initialPath?: string;
		onClose: () => void;
		onUse: (_path: string) => void;
	}

	let { initialPath, onClose, onUse }: Props = $props();
	const queryClient = createAppQueryClient();
</script>

<QueryClientProvider client={queryClient}>
	<LibraryFolderPicker {initialPath} {onClose} {onUse} />
</QueryClientProvider>
