<script lang="ts">
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import * as Table from '$lib/components/ui/table';
	import ConfirmActionButton from '$lib/components/shared/ConfirmActionButton.svelte';
	import type { DLNARendererProfile } from '$lib/settings/types';

	interface Props {
		profiles: DLNARendererProfile[];
		search: string;
		selectedId?: string;
		onSearch: (_value: string) => void;
		onSelect: (_profile: DLNARendererProfile) => void;
		onClone: (_profile: DLNARendererProfile) => void;
		onReset: (_profile: DLNARendererProfile) => void;
		onExport: (_profile: DLNARendererProfile) => void;
	}

	let { profiles, search, selectedId, onSearch, onSelect, onClone, onReset, onExport }: Props =
		$props();
</script>

<section class="grid gap-3" aria-label="DLNA renderer profiles">
	<div class="flex flex-wrap items-center justify-between gap-3">
		<h3 class="m-0 text-sm font-semibold">Device profiles</h3>
		<Input
			class="w-full sm:w-72"
			aria-label="Search profiles"
			placeholder="Search profiles"
			value={search}
			oninput={(event) => onSearch(event.currentTarget.value)}
		/>
	</div>
	<Table.Root>
		<Table.Header>
			<Table.Row>
				<Table.Head>Name</Table.Head>
				<Table.Head>Family</Table.Head>
				<Table.Head>Class</Table.Head>
				<Table.Head>Enabled</Table.Head>
				<Table.Head>Priority</Table.Head>
				<Table.Head>Version</Table.Head>
				<Table.Head>Customized</Table.Head>
				<Table.Head class="text-right">Actions</Table.Head>
			</Table.Row>
		</Table.Header>
		<Table.Body>
			{#each profiles as profile (profile.id)}
				<Table.Row class={selectedId === profile.id ? 'bg-muted/50' : ''}>
					<Table.Cell>
						<Button
							type="button"
							variant="link"
							class="h-auto p-0 text-left font-medium"
							onclick={() => onSelect(profile)}
						>
							{profile.name}
						</Button>
					</Table.Cell>
					<Table.Cell>{profile.vendor || 'Generic'}</Table.Cell>
					<Table.Cell>{profile.deviceClass}</Table.Cell>
					<Table.Cell>
						<Badge variant={profile.enabled ? 'default' : 'secondary'}>
							{profile.enabled ? 'Enabled' : 'Disabled'}
						</Badge>
					</Table.Cell>
					<Table.Cell>{profile.priority}</Table.Cell>
					<Table.Cell>{profile.sourceVersion}</Table.Cell>
					<Table.Cell>{profile.customized ? 'Yes' : 'No'}</Table.Cell>
					<Table.Cell>
						<div class="flex justify-end gap-2">
							<Button type="button" size="sm" variant="outline" onclick={() => onClone(profile)}>
								Clone
							</Button>
							<ConfirmActionButton
								label={`Reset ${profile.name}`}
								title="Reset profile"
								description={`Reset ${profile.name} to seeded defaults?`}
								confirmLabel="Reset"
								confirmingLabel="Resetting"
								variant="outline"
								size="sm"
								disabled={profile.source !== 'mema_seed'}
								onConfirm={() => onReset(profile)}
							>
								Reset
							</ConfirmActionButton>
							<Button type="button" size="sm" variant="outline" onclick={() => onExport(profile)}>
								Export
							</Button>
						</div>
					</Table.Cell>
				</Table.Row>
			{:else}
				<Table.Row>
					<Table.Cell colspan={8} class="py-8 text-center text-muted-foreground">
						No profiles match search
					</Table.Cell>
				</Table.Row>
			{/each}
		</Table.Body>
	</Table.Root>
</section>
