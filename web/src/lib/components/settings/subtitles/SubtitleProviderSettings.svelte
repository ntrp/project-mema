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
	const provider = $derived(providers.find((item) => item.type === 'opensubtitles'));
</script>

<div class="grid items-start gap-[18px] md:grid-cols-2">
	{#key provider?.id ?? 'opensubtitles'}
		<SubtitleProviderCard
			{provider}
			{onSave}
			{onDelete}
			{onTest}
			{testingId}
			{savingId}
			testResult={provider ? testResults[provider.id] : undefined}
		/>
	{/key}
</div>
