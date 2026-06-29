<script lang="ts">
	import MetadataProviderCard from './MetadataProviderCard.svelte';
	import type {
		IntegrationTestResults,
		MetadataProvider,
		MetadataProviderForm,
		MetadataProviderType
	} from '$lib/settings/types';

	interface ProviderDefinition {
		type: MetadataProviderType;
		name: string;
		baseUrl: string;
		priority: number;
		fields: 'tmdb' | 'tvdb';
	}

	interface Props {
		metadataProviders: MetadataProvider[];
		onSave: (_form: MetadataProviderForm) => void | Promise<void>;
		onTest: (_id: string) => void | Promise<void>;
		testingId?: string;
		savingId?: string;
		testResults: IntegrationTestResults;
	}

	let { metadataProviders, onSave, onTest, testingId, savingId, testResults }: Props = $props();

	const providerDefinitions: ProviderDefinition[] = [
		{
			type: 'tmdb',
			name: 'TMDB',
			baseUrl: 'https://api.themoviedb.org/3',
			priority: 100,
			fields: 'tmdb'
		},
		{
			type: 'tvdb',
			name: 'TVDB',
			baseUrl: 'https://api4.thetvdb.com/v4',
			priority: 110,
			fields: 'tvdb'
		}
	];
</script>

<div class="provider-grid">
	{#each providerDefinitions as definition (definition.type)}
		{@const provider = metadataProviders.find((item) => item.type === definition.type)}
		{@const key = provider?.id ?? definition.type}
		{#key key}
			<MetadataProviderCard
				{definition}
				{provider}
				{onSave}
				{onTest}
				{testingId}
				{savingId}
				testResult={provider ? testResults[provider.id] : undefined}
			/>
		{/key}
	{/each}
</div>
