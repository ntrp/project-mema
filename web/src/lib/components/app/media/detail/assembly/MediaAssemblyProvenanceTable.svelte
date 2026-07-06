<script lang="ts">
	import { Badge } from '$lib/components/ui/badge';
	import * as Table from '$lib/components/ui/table';
	import type { MediaComponentAssemblyInput } from '$lib/settings/types';
	import { fileName } from './mediaComponentAssemblyView';

	interface Props {
		inputs: MediaComponentAssemblyInput[];
	}

	let { inputs }: Props = $props();

	function provenanceLabel(input: MediaComponentAssemblyInput) {
		const source = input.provenance.sourceFilePath;
		const stream = input.provenance.streamId;
		const language = input.provenance.language;
		return [
			typeof source === 'string' ? fileName(source) : undefined,
			streamLabel(stream),
			language
		]
			.filter(Boolean)
			.join(' · ');
	}

	function streamLabel(stream: unknown) {
		return typeof stream === 'number' ? `stream ${stream}` : undefined;
	}
</script>

<div class="overflow-hidden rounded-md border">
	<Table.Root>
		<Table.Header>
			<Table.Row>
				<Table.Head>Stream</Table.Head>
				<Table.Head>Input</Table.Head>
				<Table.Head>Provenance</Table.Head>
			</Table.Row>
		</Table.Header>
		<Table.Body>
			{#each inputs as input (input.id)}
				<Table.Row>
					<Table.Cell
						><Badge variant={input.streamType === 'video' ? 'default' : 'secondary'}
							>{input.streamType}</Badge
						></Table.Cell
					>
					<Table.Cell class="max-w-[280px] truncate">{fileName(input.inputPath)}</Table.Cell>
					<Table.Cell class="text-muted-foreground"
						>{provenanceLabel(input) || 'Recorded'}</Table.Cell
					>
				</Table.Row>
			{/each}
		</Table.Body>
	</Table.Root>
</div>
