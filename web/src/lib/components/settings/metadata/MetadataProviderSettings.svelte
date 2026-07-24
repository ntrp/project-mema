<script lang="ts">
	import MetadataProviderCard from './MetadataProviderCard.svelte';
	import { metadataProviderDefinitions } from './metadataProviderDefinitions';
	import type {
		IntegrationTestResults,
		MetadataProvider,
		MetadataProviderForm
	} from '$lib/settings/types';

	interface Props {
		metadataProviders: MetadataProvider[];
		onSave: (_form: MetadataProviderForm) => void | Promise<void>;
		onTest: (_id: string) => void | Promise<void>;
		testingId?: string;
		savingId?: string;
		testResults: IntegrationTestResults;
	}

	let { metadataProviders, onSave, onTest, testingId, savingId, testResults }: Props = $props();
</script>

<div class="grid items-start gap-[18px] md:grid-cols-2">
	{#each metadataProviderDefinitions as definition (definition.type)}
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
