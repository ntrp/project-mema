<script lang="ts">
	import { resolve } from '$app/paths';
	import SettingsSelect from '$lib/components/settings/shared/SettingsSelect.svelte';
	import EmptyState from '$lib/components/shared/EmptyState.svelte';
	import PageHeading from '$lib/components/shared/PageHeading.svelte';
	import SectionHeading from '$lib/components/shared/SectionHeading.svelte';
	import StatusPill from '$lib/components/shared/StatusPill.svelte';
	import { Button } from '$lib/components/ui/button';
	import MediaRequestCard from './MediaRequestCard.svelte';
	import type {
		LibraryFolder,
		MediaRequest,
		MediaRequestApproveRequest,
		QualityProfileOption
	} from '$lib/settings/types';

	interface Props {
		requests: MediaRequest[];
		selectedRequestId?: string;
		libraryFolders: LibraryFolder[];
		qualityProfiles: QualityProfileOption[];
		canManage: boolean;
		approvingRequestId?: string;
		onApprove: (_request: MediaRequest, _approval: MediaRequestApproveRequest) => void;
	}

	let {
		requests,
		selectedRequestId,
		libraryFolders,
		qualityProfiles,
		canManage,
		approvingRequestId,
		onApprove
	}: Props = $props();

	const selectedRequest = $derived(
		selectedRequestId ? requests.find((request) => request.id === selectedRequestId) : undefined
	);

	let qualityProfileId = $state('');
	let libraryFolderId = $state('');
	let qualityProfileOptions = $derived([
		{ value: '', label: 'Select profile' },
		...qualityProfiles.map((profile) => ({ value: profile.id, label: profile.name }))
	]);
	let libraryFolderOptions = $derived([
		{ value: '', label: 'Select folder' },
		...libraryFolders.map((folder) => ({ value: folder.id, label: folder.path }))
	]);
	let requestFacts = $derived(
		selectedRequest
			? [
					{ label: 'Requested by', value: selectedRequest.requestedByUsername },
					{ label: 'Year', value: selectedRequest.year ?? 'Unknown' },
					{ label: 'Quality profile', value: profileName(selectedRequest.qualityProfileId) },
					{ label: 'Library folder', value: folderName(selectedRequest.libraryFolderId) }
				]
			: []
	);

	function folderName(id?: string) {
		if (!id) {
			return 'Not selected';
		}
		return libraryFolders.find((folder) => folder.id === id)?.path ?? id;
	}

	function profileName(id?: string) {
		if (!id) {
			return 'Not selected';
		}
		return qualityProfiles.find((profile) => profile.id === id)?.name ?? id;
	}

	function approve(event: SubmitEvent) {
		event.preventDefault();
		if (!selectedRequest || !qualityProfileId || !libraryFolderId) {
			return;
		}
		onApprove(selectedRequest, { qualityProfileId, libraryFolderId });
	}
</script>

{#if selectedRequestId}
	<a
		class="w-fit font-extrabold text-primary no-underline hover:underline focus-visible:underline focus-visible:outline-none"
		href={resolve('/requests')}>Back to requests</a
	>
	{#if selectedRequest}
		<section
			class="grid gap-[18px] rounded-md border border-border bg-card p-5"
			aria-label={selectedRequest.title}
		>
			<SectionHeading title={selectedRequest.title} kicker={selectedRequest.status}>
				{#snippet actions()}
					<StatusPill>{selectedRequest.type}</StatusPill>
				{/snippet}
			</SectionHeading>
			<div class="grid gap-3 md:grid-cols-4">
				{#each requestFacts as fact (fact.label)}
					<div class="grid gap-1 rounded-md border border-border bg-muted p-2.5">
						<strong>{fact.label}</strong>
						<span class="break-words text-muted-foreground">{fact.value}</span>
					</div>
				{/each}
			</div>
			{#if selectedRequest.overview}
				<p>{selectedRequest.overview}</p>
			{/if}
			{#if selectedRequest.tags?.length}
				<div class="flex flex-wrap gap-2" aria-label="Tags">
					{#each selectedRequest.tags as tag (tag)}
						<StatusPill>{tag}</StatusPill>
					{/each}
				</div>
			{/if}

			{#if canManage && selectedRequest.status === 'pending'}
				<form class="grid gap-4 md:grid-cols-2" onsubmit={approve}>
					<label class="grid gap-1.5">
						<span class="text-sm font-bold text-muted-foreground">Quality profile</span>
						<SettingsSelect
							value={qualityProfileId}
							options={qualityProfileOptions}
							onValueChange={(value) => (qualityProfileId = value)}
						/>
					</label>
					<label class="grid gap-1.5">
						<span class="text-sm font-bold text-muted-foreground">Library folder</span>
						<SettingsSelect
							value={libraryFolderId}
							options={libraryFolderOptions}
							onValueChange={(value) => (libraryFolderId = value)}
						/>
					</label>
					<div class="flex items-center gap-3 md:col-span-2">
						<Button
							type="submit"
							disabled={approvingRequestId === selectedRequest.id || libraryFolders.length === 0}
						>
							{approvingRequestId === selectedRequest.id ? 'Approving' : 'Approve'}
						</Button>
					</div>
				</form>
			{/if}
		</section>
	{:else}
		<EmptyState
			title="Request not found"
			description="The request is not visible to your account."
		/>
	{/if}
{:else}
	<PageHeading eyebrow="Requests" title="Media requests" titleId="home-title" />
	{#if requests.length > 0}
		<div class="grid gap-2.5">
			{#each requests as request (request.id)}
				<MediaRequestCard {request} />
			{/each}
		</div>
	{:else}
		<EmptyState title="No requests" description="Requested media will appear here." />
	{/if}
{/if}
