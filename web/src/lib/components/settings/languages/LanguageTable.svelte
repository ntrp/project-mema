<script lang="ts">
	import SettingsRowActionButton from '$lib/components/settings/shared/SettingsRowActionButton.svelte';
	import { Badge } from '$lib/components/ui/badge';
	import { Card } from '$lib/components/ui/card';
	import * as Table from '$lib/components/ui/table';
	import { formatDate } from '$lib/settings/dateFormat';
	import type { Language } from '$lib/settings/types';

	interface Props {
		languages: Language[];
		deletingCode?: string;
		onEdit: (_language: Language) => void;
		onDelete: (_language: Language) => void;
	}

	let { languages, deletingCode, onEdit, onDelete }: Props = $props();
</script>

<Card class="gap-0 p-0" aria-label="Languages">
	<Table.Root>
		<Table.Header>
			<Table.Row>
				<Table.Head>Code</Table.Head>
				<Table.Head>Display name</Table.Head>
				<Table.Head>Aliases</Table.Head>
				<Table.Head>Updated</Table.Head>
				<Table.Head class="text-right">Actions</Table.Head>
			</Table.Row>
		</Table.Header>
		<Table.Body>
			{#each languages as language (language.code)}
				<Table.Row>
					<Table.Cell><Badge variant="secondary">{language.code}</Badge></Table.Cell>
					<Table.Cell>{language.displayName}</Table.Cell>
					<Table.Cell class="max-w-[420px] truncate text-muted-foreground">
						{language.aliases.join(', ')}
					</Table.Cell>
					<Table.Cell>{formatDate(language.updatedAt)}</Table.Cell>
					<Table.Cell>
						<div class="flex justify-end gap-2">
							<SettingsRowActionButton
								label={`Edit ${language.displayName}`}
								icon="edit"
								onclick={() => onEdit(language)}
							/>
							<SettingsRowActionButton
								label={`${deletingCode === language.code ? 'Deleting' : 'Delete'} ${language.displayName}`}
								icon="delete"
								variant="destructive"
								disabled={deletingCode === language.code}
								onclick={() => onDelete(language)}
							/>
						</div>
					</Table.Cell>
				</Table.Row>
			{:else}
				<Table.Row>
					<Table.Cell colspan={5} class="py-8 text-center text-muted-foreground">
						No languages configured
					</Table.Cell>
				</Table.Row>
			{/each}
		</Table.Body>
	</Table.Root>
</Card>
