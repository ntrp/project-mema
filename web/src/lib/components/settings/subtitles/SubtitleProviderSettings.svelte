<script lang="ts">
	import SubtitleProviderCard from './SubtitleProviderCard.svelte';
	import type {
		IntegrationTestResults,
		SubtitleProvider,
		SubtitleProviderForm
	} from '$lib/settings/types';

	interface Props {
		providers: SubtitleProvider[];
		onSave: (_form: SubtitleProviderForm) => void | Promise<void>;
		onDelete: (_id: string) => void | Promise<void>;
		onTest: (_id: string) => void | Promise<void>;
		testingId?: string;
		savingId?: string;
		testResults: IntegrationTestResults;
	}

	let { providers, onSave, onDelete, onTest, testingId, savingId, testResults }: Props = $props();
	const openSubtitlesProvider = $derived(providers.find((item) => item.type === 'opensubtitles'));
	const mockProvider = $derived(providers.find((item) => item.type === 'mock'));
</script>

<div class="grid items-start gap-[18px] md:grid-cols-2">
	{#key openSubtitlesProvider?.id ?? 'opensubtitles'}
		<SubtitleProviderCard
			provider={openSubtitlesProvider}
			providerType="opensubtitles"
			{onSave}
			{onDelete}
			{onTest}
			{testingId}
			{savingId}
			testResult={openSubtitlesProvider ? testResults[openSubtitlesProvider.id] : undefined}
		/>
	{/key}
	{#key mockProvider?.id ?? 'mock'}
		<SubtitleProviderCard
			provider={mockProvider}
			providerType="mock"
			{onSave}
			{onDelete}
			{onTest}
			{testingId}
			{savingId}
			testResult={mockProvider ? testResults[mockProvider.id] : undefined}
		/>
	{/key}
</div>
