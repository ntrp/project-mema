<script lang="ts">
	import NoticeStack from '$lib/components/settings/shared/NoticeStack.svelte';
	import { Button } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
	import * as Table from '$lib/components/ui/table';
	import { createQualitySizeResources } from './resources.svelte';
	import { groupQualitiesByResolution } from '$lib/settings/qualityGroups';
	import type { QualitySizeSetting } from '$lib/settings/types';
	import QualitySizeRow from './QualitySizeRow.svelte';
	import { qualityRequest, rowError } from './qualitySize';

	const resources = createQualitySizeResources();
	let qualities = $derived<QualitySizeSetting[]>(resources.query.data ?? []);
	const loading = $derived(resources.query.isFetching);
	const saving = $derived(resources.update.isPending);
	let errorMessage = $state('');
	let message = $state('');

	const hasValidationErrors = $derived(qualities.some((quality) => rowError(quality) !== ''));
	const qualityGroups = $derived(groupQualitiesByResolution(qualities));

	async function loadQualitySizes() {
		errorMessage = '';
		try {
			await resources.query.refetch();
		} catch (error) {
			errorMessage = error instanceof Error ? error.message : 'Could not load quality sizes';
		}
	}

	async function saveQualitySizes(event: SubmitEvent) {
		event.preventDefault();
		message = '';
		errorMessage = '';
		if (hasValidationErrors) {
			errorMessage = 'Fix invalid quality sizes before saving';
			return;
		}

		try {
			const response = await resources.update.mutateAsync(qualities.map(qualityRequest));
			qualities = response.qualities;
			message = 'Quality sizes saved';
		} catch (error) {
			errorMessage = error instanceof Error ? error.message : 'Could not save quality sizes';
		}
	}

	function updateQuality(nextQuality: QualitySizeSetting) {
		qualities = qualities.map((quality) =>
			quality.qualityId === nextQuality.qualityId ? nextQuality : quality
		);
		message = '';
	}
</script>

<Card.Root aria-labelledby="quality-size-title">
	<form class="grid gap-3.5" onsubmit={saveQualitySizes}>
		<Card.Header>
			<div>
				<Card.Description>Release scoring</Card.Description>
				<Card.Title id="quality-size-title">Quality sizes</Card.Title>
			</div>
			<Card.Action>
				<div class="flex flex-wrap justify-end gap-2">
					<Button
						type="button"
						variant="outline"
						disabled={loading || saving}
						onclick={loadQualitySizes}
					>
						Reload
					</Button>
					<Button type="submit" disabled={loading || saving || hasValidationErrors}>
						{saving ? 'Saving' : 'Save sizes'}
					</Button>
				</div>
			</Card.Action>
		</Card.Header>

		<Card.Content class="grid gap-4">
			<NoticeStack {message} {errorMessage} />

			<Table.Root class="min-w-[980px]">
				<Table.Header>
					<Table.Row>
						<Table.Head class="w-[180px] min-w-40">Quality</Table.Head>
						<Table.Head>Size limit</Table.Head>
					</Table.Row>
				</Table.Header>
				<Table.Body>
					{#if loading}
						<Table.Row>
							<Table.Cell colspan={2} class="text-muted-foreground"
								>Loading quality sizes</Table.Cell
							>
						</Table.Row>
					{:else if qualities.length === 0}
						<Table.Row>
							<Table.Cell colspan={2} class="text-muted-foreground">No qualities loaded</Table.Cell>
						</Table.Row>
					{:else}
						{#each qualityGroups as group (group.id)}
							<Table.Row>
								<Table.Head
									colspan={2}
									class="bg-muted/60 text-xs font-black text-muted-foreground uppercase"
								>
									{group.label}
								</Table.Head>
							</Table.Row>
							{#each group.qualities as quality (quality.qualityId)}
								<QualitySizeRow {quality} onChange={updateQuality} />
							{/each}
						{/each}
					{/if}
				</Table.Body>
			</Table.Root>
		</Card.Content>
	</form>
</Card.Root>
