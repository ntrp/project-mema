<script lang="ts">
	import { resolve } from '$app/paths';
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

	function posterUrl(path?: string) {
		if (!path) {
			return undefined;
		}
		if (path.startsWith('http://') || path.startsWith('https://')) {
			return path;
		}
		return `https://image.tmdb.org/t/p/w185${path}`;
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
	<a class="back-link" href={resolve('/requests')}>Back to requests</a>
	{#if selectedRequest}
		<section class="media-request-detail panel" aria-labelledby="request-detail-title">
			<div class="section-heading">
				<div>
					<p class="section-kicker">{selectedRequest.status}</p>
					<h1 id="request-detail-title">{selectedRequest.title}</h1>
				</div>
				<span class="status-pill">{selectedRequest.type}</span>
			</div>
			<div class="request-detail-grid">
				<div>
					<strong>Requested by</strong>
					<span>{selectedRequest.requestedByUsername}</span>
				</div>
				<div>
					<strong>Year</strong>
					<span>{selectedRequest.year ?? 'Unknown'}</span>
				</div>
				<div>
					<strong>Quality profile</strong>
					<span>{profileName(selectedRequest.qualityProfileId)}</span>
				</div>
				<div>
					<strong>Library folder</strong>
					<span>{folderName(selectedRequest.libraryFolderId)}</span>
				</div>
			</div>
			{#if selectedRequest.overview}
				<p>{selectedRequest.overview}</p>
			{/if}
			{#if selectedRequest.tags?.length}
				<div class="media-tags" aria-label="Tags">
					{#each selectedRequest.tags as tag (tag)}
						<span>{tag}</span>
					{/each}
				</div>
			{/if}

			{#if canManage && selectedRequest.status === 'pending'}
				<form class="settings-form compact-form" onsubmit={approve}>
					<label>
						<span>Quality profile</span>
						<select bind:value={qualityProfileId}>
							<option value="" disabled>Select profile</option>
							{#each qualityProfiles as profile (profile.id)}
								<option value={profile.id}>{profile.name}</option>
							{/each}
						</select>
					</label>
					<label>
						<span>Library folder</span>
						<select bind:value={libraryFolderId}>
							<option value="" disabled>Select folder</option>
							{#each libraryFolders as folder (folder.id)}
								<option value={folder.id}>{folder.path}</option>
							{/each}
						</select>
					</label>
					<div class="form-actions wide">
						<button
							type="submit"
							disabled={approvingRequestId === selectedRequest.id || libraryFolders.length === 0}
						>
							{approvingRequestId === selectedRequest.id ? 'Approving' : 'Approve'}
						</button>
					</div>
				</form>
			{/if}
		</section>
	{:else}
		<section class="empty-state">
			<h2>Request not found</h2>
			<p>The request is not visible to your account.</p>
		</section>
	{/if}
{:else}
	<div class="page-heading">
		<p>Requests</p>
		<h1 id="home-title">Media requests</h1>
	</div>
	{#if requests.length > 0}
		<div class="wide-card-list">
			{#each requests as request (request.id)}
				<a
					class="wide-media-card request-card"
					href={resolve('/requests/[id]', { id: request.id })}
				>
					<div class="wide-poster">
						{#if posterUrl(request.posterPath)}
							<img src={posterUrl(request.posterPath)} alt="" loading="lazy" />
						{:else}
							<div class="poster-placeholder">{request.type}</div>
						{/if}
					</div>
					<div class="wide-media-body">
						<h3>{request.title}</h3>
						<p>
							{request.type}{request.year ? ` · ${request.year}` : ''} · Requested by {request.requestedByUsername}
						</p>
						{#if request.overview}
							<p>{request.overview}</p>
						{/if}
						{#if request.tags?.length}
							<div class="media-tags compact-tags" aria-label="Tags">
								{#each request.tags.slice(0, 3) as tag (tag)}
									<span>{tag}</span>
								{/each}
							</div>
						{/if}
					</div>
					<span class="status-pill">{request.status}</span>
				</a>
			{/each}
		</div>
	{:else}
		<section class="empty-state">
			<h2>No requests</h2>
			<p>Requested media will appear here.</p>
		</section>
	{/if}
{/if}
