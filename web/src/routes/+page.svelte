<script lang="ts">
	import { onMount } from 'svelte';

	import { client } from '$lib/api/client';
	import type { components } from '$lib/api/generated/schema';

	type DownloadClient = components['schemas']['DownloadClient'];
	type DownloadClientRequest = components['schemas']['DownloadClientRequest'];
	type DownloadClientType = components['schemas']['DownloadClientType'];
	type Indexer = components['schemas']['Indexer'];
	type IndexerRequest = components['schemas']['IndexerRequest'];
	type IndexerType = components['schemas']['IndexerType'];

	type DownloadClientForm = DownloadClientRequest & { id?: string };
	type IndexerForm = Omit<IndexerRequest, 'categories'> & { id?: string; categoriesText: string };

	let authenticated = $state(false);
	let loading = $state(true);
	let savingDownloadClient = $state(false);
	let savingIndexer = $state(false);
	let message = $state('');
	let errorMessage = $state('');
	let username = $state('admin');
	let password = $state('admin');
	let downloadClients = $state<DownloadClient[]>([]);
	let indexers = $state<Indexer[]>([]);
	let downloadForm = $state<DownloadClientForm>(emptyDownloadClientForm());
	let indexerForm = $state<IndexerForm>(emptyIndexerForm());

	const downloadClientTypes: DownloadClientType[] = ['transmission', 'sabnzbd'];
	const indexerTypes: IndexerType[] = ['torznab', 'newznab', 'rss'];

	const enabledDownloadClients = $derived(downloadClients.filter((item) => item.enabled).length);
	const enabledIndexers = $derived(indexers.filter((item) => item.enabled).length);

	onMount(() => {
		void initialise();
	});

	async function initialise() {
		loading = true;
		errorMessage = '';

		const { data } = await client.GET('/auth/session');
		authenticated = Boolean(data?.authenticated);

		if (authenticated) {
			await loadSettings();
		}

		loading = false;
	}

	async function login(event: SubmitEvent) {
		event.preventDefault();
		errorMessage = '';
		message = '';

		const { data, error } = await client.POST('/auth/login', {
			body: { username, password }
		});

		if (error || !data?.authenticated) {
			errorMessage = error?.message ?? 'Login failed';
			return;
		}

		authenticated = true;
		await loadSettings();
	}

	async function loadSettings() {
		const [clientResult, indexerResult] = await Promise.all([
			client.GET('/settings/download-clients'),
			client.GET('/settings/indexers')
		]);

		if (clientResult.error) {
			errorMessage = clientResult.error.message;
			return;
		}
		if (indexerResult.error) {
			errorMessage = indexerResult.error.message;
			return;
		}

		downloadClients = clientResult.data?.clients ?? [];
		indexers = indexerResult.data?.indexers ?? [];
	}

	async function saveDownloadClient(event: SubmitEvent) {
		event.preventDefault();
		savingDownloadClient = true;
		errorMessage = '';
		message = '';

		const body = normalizeDownloadClientForm(downloadForm);
		const result = downloadForm.id
			? await client.PUT('/settings/download-clients/{id}', {
					params: { path: { id: downloadForm.id } },
					body
				})
			: await client.POST('/settings/download-clients', { body });

		savingDownloadClient = false;
		if (result.error) {
			errorMessage = result.error.message;
			return;
		}

		downloadForm = emptyDownloadClientForm();
		message = 'Download client saved';
		await loadSettings();
	}

	async function saveIndexer(event: SubmitEvent) {
		event.preventDefault();
		savingIndexer = true;
		errorMessage = '';
		message = '';

		const body = normalizeIndexerForm(indexerForm);
		const result = indexerForm.id
			? await client.PUT('/settings/indexers/{id}', {
					params: { path: { id: indexerForm.id } },
					body
				})
			: await client.POST('/settings/indexers', { body });

		savingIndexer = false;
		if (result.error) {
			errorMessage = result.error.message;
			return;
		}

		indexerForm = emptyIndexerForm();
		message = 'Indexer saved';
		await loadSettings();
	}

	async function deleteDownloadClient(id: string) {
		errorMessage = '';
		message = '';
		const { error } = await client.DELETE('/settings/download-clients/{id}', {
			params: { path: { id } }
		});

		if (error) {
			errorMessage = error.message;
			return;
		}

		if (downloadForm.id === id) {
			downloadForm = emptyDownloadClientForm();
		}
		message = 'Download client deleted';
		await loadSettings();
	}

	async function deleteIndexer(id: string) {
		errorMessage = '';
		message = '';
		const { error } = await client.DELETE('/settings/indexers/{id}', {
			params: { path: { id } }
		});

		if (error) {
			errorMessage = error.message;
			return;
		}

		if (indexerForm.id === id) {
			indexerForm = emptyIndexerForm();
		}
		message = 'Indexer deleted';
		await loadSettings();
	}

	function editDownloadClient(clientItem: DownloadClient) {
		downloadForm = {
			id: clientItem.id,
			name: clientItem.name,
			type: clientItem.type,
			baseUrl: clientItem.baseUrl,
			username: clientItem.username ?? '',
			password: clientItem.password ?? '',
			apiKey: clientItem.apiKey ?? '',
			category: clientItem.category ?? '',
			enabled: clientItem.enabled,
			priority: clientItem.priority
		};
	}

	function editIndexer(indexer: Indexer) {
		indexerForm = {
			id: indexer.id,
			name: indexer.name,
			type: indexer.type,
			baseUrl: indexer.baseUrl,
			apiKey: indexer.apiKey ?? '',
			categoriesText: (indexer.categories ?? []).join(', '),
			enabled: indexer.enabled,
			priority: indexer.priority
		};
	}

	function emptyDownloadClientForm(): DownloadClientForm {
		return {
			name: '',
			type: 'transmission',
			baseUrl: '',
			username: '',
			password: '',
			apiKey: '',
			category: '',
			enabled: true,
			priority: 100
		};
	}

	function emptyIndexerForm(): IndexerForm {
		return {
			name: '',
			type: 'torznab',
			baseUrl: '',
			apiKey: '',
			categoriesText: '',
			enabled: true,
			priority: 100
		};
	}

	function normalizeDownloadClientForm(form: DownloadClientForm): DownloadClientRequest {
		return {
			name: form.name.trim(),
			type: form.type,
			baseUrl: form.baseUrl.trim(),
			username: optionalString(form.username),
			password: optionalString(form.password),
			apiKey: optionalString(form.apiKey),
			category: optionalString(form.category),
			enabled: form.enabled,
			priority: form.priority
		};
	}

	function normalizeIndexerForm(form: IndexerForm): IndexerRequest {
		return {
			name: form.name.trim(),
			type: form.type,
			baseUrl: form.baseUrl.trim(),
			apiKey: optionalString(form.apiKey),
			categories: parseCategories(form.categoriesText),
			enabled: form.enabled,
			priority: form.priority
		};
	}

	function optionalString(value: string | undefined) {
		const trimmed = value?.trim() ?? '';
		return trimmed === '' ? undefined : trimmed;
	}

	function parseCategories(value: string) {
		return value
			.split(',')
			.map((item) => Number.parseInt(item.trim(), 10))
			.filter((item) => Number.isInteger(item));
	}
</script>

<svelte:head>
	<title>Media Manager Settings</title>
	<meta
		name="description"
		content="Configure download clients and indexers for the self-hosted media manager"
	/>
</svelte:head>

<main class="shell">
	<header>
		<div>
			<p>Settings</p>
			<h1>Download clients and indexers</h1>
		</div>
		<div class="summary">
			<span>{downloadClients.length} clients</span>
			<span>{enabledDownloadClients} enabled</span>
			<span>{indexers.length} indexers</span>
			<span>{enabledIndexers} enabled</span>
		</div>
	</header>

	{#if loading}
		<section class="panel">
			<p class="muted">Loading settings</p>
		</section>
	{:else if !authenticated}
		<section class="panel auth-panel" aria-labelledby="login-title">
			<h2 id="login-title">Admin login</h2>
			<form onsubmit={login}>
				<label>
					<span>Username</span>
					<input bind:value={username} autocomplete="username" required />
				</label>
				<label>
					<span>Password</span>
					<input bind:value={password} autocomplete="current-password" type="password" required />
				</label>
				<button type="submit">Log in</button>
			</form>
		</section>
	{:else}
		{#if message}
			<p class="notice success">{message}</p>
		{/if}
		{#if errorMessage}
			<p class="notice error">{errorMessage}</p>
		{/if}

		<section class="settings-grid">
			<div class="panel" aria-labelledby="download-client-form-title">
				<div class="section-heading">
					<h2 id="download-client-form-title">
						{downloadForm.id ? 'Edit download client' : 'Add download client'}
					</h2>
					{#if downloadForm.id}
						<button
							type="button"
							class="secondary"
							onclick={() => (downloadForm = emptyDownloadClientForm())}
						>
							Cancel
						</button>
					{/if}
				</div>

				<form class="settings-form" onsubmit={saveDownloadClient}>
					<label>
						<span>Name</span>
						<input bind:value={downloadForm.name} required maxlength="200" />
					</label>
					<label>
						<span>Type</span>
						<select bind:value={downloadForm.type}>
							{#each downloadClientTypes as type (type)}
								<option value={type}>{type}</option>
							{/each}
						</select>
					</label>
					<label class="wide">
						<span>Base URL</span>
						<input bind:value={downloadForm.baseUrl} placeholder="http://host:port" required />
					</label>
					<label>
						<span>Username</span>
						<input bind:value={downloadForm.username} autocomplete="off" />
					</label>
					<label>
						<span>Password</span>
						<input bind:value={downloadForm.password} autocomplete="off" type="password" />
					</label>
					<label>
						<span>API key</span>
						<input bind:value={downloadForm.apiKey} autocomplete="off" />
					</label>
					<label>
						<span>Category</span>
						<input bind:value={downloadForm.category} placeholder="movies" />
					</label>
					<label>
						<span>Priority</span>
						<input bind:value={downloadForm.priority} min="0" max="1000" type="number" />
					</label>
					<label class="toggle">
						<input bind:checked={downloadForm.enabled} type="checkbox" />
						<span>Enabled</span>
					</label>
					<button type="submit" disabled={savingDownloadClient}>
						{savingDownloadClient ? 'Saving' : 'Save client'}
					</button>
				</form>
			</div>

			<div class="panel" aria-labelledby="indexer-form-title">
				<div class="section-heading">
					<h2 id="indexer-form-title">{indexerForm.id ? 'Edit indexer' : 'Add indexer'}</h2>
					{#if indexerForm.id}
						<button
							type="button"
							class="secondary"
							onclick={() => (indexerForm = emptyIndexerForm())}
						>
							Cancel
						</button>
					{/if}
				</div>

				<form class="settings-form" onsubmit={saveIndexer}>
					<label>
						<span>Name</span>
						<input bind:value={indexerForm.name} required maxlength="200" />
					</label>
					<label>
						<span>Type</span>
						<select bind:value={indexerForm.type}>
							{#each indexerTypes as type (type)}
								<option value={type}>{type}</option>
							{/each}
						</select>
					</label>
					<label class="wide">
						<span>Base URL</span>
						<input
							bind:value={indexerForm.baseUrl}
							placeholder="https://indexer.example"
							required
						/>
					</label>
					<label>
						<span>API key</span>
						<input bind:value={indexerForm.apiKey} autocomplete="off" />
					</label>
					<label>
						<span>Categories</span>
						<input bind:value={indexerForm.categoriesText} placeholder="2000, 5000" />
					</label>
					<label>
						<span>Priority</span>
						<input bind:value={indexerForm.priority} min="0" max="1000" type="number" />
					</label>
					<label class="toggle">
						<input bind:checked={indexerForm.enabled} type="checkbox" />
						<span>Enabled</span>
					</label>
					<button type="submit" disabled={savingIndexer}>
						{savingIndexer ? 'Saving' : 'Save indexer'}
					</button>
				</form>
			</div>
		</section>

		<section class="list-grid">
			<div class="panel" aria-labelledby="download-client-list-title">
				<h2 id="download-client-list-title">Download clients</h2>
				<div class="table-wrap">
					<table>
						<thead>
							<tr>
								<th>Name</th>
								<th>Type</th>
								<th>Base URL</th>
								<th>Priority</th>
								<th>Status</th>
								<th></th>
							</tr>
						</thead>
						<tbody>
							{#each downloadClients as item (item.id)}
								<tr>
									<td>{item.name}</td>
									<td>{item.type}</td>
									<td>{item.baseUrl}</td>
									<td>{item.priority}</td>
									<td>{item.enabled ? 'Enabled' : 'Disabled'}</td>
									<td class="row-actions">
										<button
											type="button"
											class="secondary"
											onclick={() => editDownloadClient(item)}
										>
											Edit
										</button>
										<button
											type="button"
											class="danger"
											onclick={() => deleteDownloadClient(item.id)}
										>
											Delete
										</button>
									</td>
								</tr>
							{:else}
								<tr>
									<td colspan="6" class="empty">No download clients configured</td>
								</tr>
							{/each}
						</tbody>
					</table>
				</div>
			</div>

			<div class="panel" aria-labelledby="indexer-list-title">
				<h2 id="indexer-list-title">Indexers</h2>
				<div class="table-wrap">
					<table>
						<thead>
							<tr>
								<th>Name</th>
								<th>Type</th>
								<th>Base URL</th>
								<th>Categories</th>
								<th>Priority</th>
								<th>Status</th>
								<th></th>
							</tr>
						</thead>
						<tbody>
							{#each indexers as item (item.id)}
								<tr>
									<td>{item.name}</td>
									<td>{item.type}</td>
									<td>{item.baseUrl}</td>
									<td>{(item.categories ?? []).join(', ') || '-'}</td>
									<td>{item.priority}</td>
									<td>{item.enabled ? 'Enabled' : 'Disabled'}</td>
									<td class="row-actions">
										<button type="button" class="secondary" onclick={() => editIndexer(item)}
											>Edit</button
										>
										<button type="button" class="danger" onclick={() => deleteIndexer(item.id)}>
											Delete
										</button>
									</td>
								</tr>
							{:else}
								<tr>
									<td colspan="7" class="empty">No indexers configured</td>
								</tr>
							{/each}
						</tbody>
					</table>
				</div>
			</div>
		</section>
	{/if}
</main>

<style>
	:global(body) {
		margin: 0;
		font-family:
			Inter,
			ui-sans-serif,
			system-ui,
			-apple-system,
			BlinkMacSystemFont,
			'Segoe UI',
			sans-serif;
		background: #f4f6f8;
		color: #151922;
	}

	.shell {
		min-height: 100vh;
		padding: 32px;
		box-sizing: border-box;
	}

	header {
		display: flex;
		align-items: flex-end;
		justify-content: space-between;
		gap: 24px;
		margin-bottom: 24px;
	}

	header p {
		margin: 0 0 8px;
		font-size: 13px;
		font-weight: 700;
		text-transform: uppercase;
		color: #586475;
	}

	h1,
	h2 {
		margin: 0;
	}

	h1 {
		font-size: 34px;
		line-height: 1.12;
	}

	h2 {
		font-size: 18px;
	}

	.summary {
		display: flex;
		flex-wrap: wrap;
		gap: 8px;
		justify-content: flex-end;
	}

	.summary span {
		border: 1px solid #d8dee8;
		border-radius: 6px;
		background: #ffffff;
		padding: 7px 10px;
		font-size: 13px;
		font-weight: 700;
		color: #445064;
	}

	.panel {
		border: 1px solid #d8dee8;
		border-radius: 8px;
		background: #ffffff;
		padding: 20px;
		box-sizing: border-box;
	}

	.auth-panel {
		max-width: 420px;
	}

	.settings-grid,
	.list-grid {
		display: grid;
		grid-template-columns: repeat(2, minmax(0, 1fr));
		gap: 18px;
		margin-bottom: 18px;
	}

	.list-grid {
		align-items: start;
	}

	.section-heading {
		display: flex;
		align-items: center;
		justify-content: space-between;
		gap: 12px;
		margin-bottom: 18px;
	}

	form,
	.settings-form {
		display: grid;
		gap: 14px;
	}

	.settings-form {
		grid-template-columns: repeat(2, minmax(0, 1fr));
	}

	label {
		display: grid;
		gap: 6px;
	}

	label span {
		font-size: 13px;
		font-weight: 700;
		color: #465266;
	}

	input,
	select {
		width: 100%;
		min-height: 38px;
		border: 1px solid #cbd3df;
		border-radius: 6px;
		background: #ffffff;
		padding: 7px 10px;
		box-sizing: border-box;
		font: inherit;
	}

	input:focus,
	select:focus {
		outline: 2px solid #7aa7d9;
		outline-offset: 1px;
	}

	.wide {
		grid-column: 1 / -1;
	}

	.toggle {
		display: flex;
		align-items: center;
		gap: 8px;
	}

	.toggle input {
		width: 18px;
		min-height: 18px;
	}

	button {
		min-height: 38px;
		border: 1px solid #202938;
		border-radius: 6px;
		background: #202938;
		color: #ffffff;
		padding: 0 14px;
		font: inherit;
		font-weight: 700;
		cursor: pointer;
	}

	button:disabled {
		cursor: wait;
		opacity: 0.65;
	}

	.secondary {
		border-color: #cbd3df;
		background: #ffffff;
		color: #202938;
	}

	.danger {
		border-color: #b42318;
		background: #ffffff;
		color: #b42318;
	}

	.table-wrap {
		overflow-x: auto;
		margin-top: 16px;
	}

	table {
		width: 100%;
		border-collapse: collapse;
		font-size: 14px;
	}

	th,
	td {
		border-bottom: 1px solid #e4e8ef;
		padding: 10px 8px;
		text-align: left;
		vertical-align: top;
	}

	th {
		font-size: 12px;
		text-transform: uppercase;
		color: #617087;
	}

	td {
		word-break: break-word;
	}

	.row-actions {
		display: flex;
		gap: 8px;
		justify-content: flex-end;
		white-space: nowrap;
	}

	.notice {
		border-radius: 6px;
		padding: 10px 12px;
		font-weight: 700;
	}

	.success {
		border: 1px solid #8fc6a4;
		background: #eef8f1;
		color: #1f6b3a;
	}

	.error {
		border: 1px solid #f0a49d;
		background: #fff1f0;
		color: #a32018;
	}

	.muted,
	.empty {
		color: #687386;
	}

	@media (max-width: 980px) {
		header,
		.settings-grid,
		.list-grid {
			grid-template-columns: 1fr;
		}

		header {
			display: grid;
		}

		.summary {
			justify-content: flex-start;
		}
	}

	@media (max-width: 640px) {
		.shell {
			padding: 18px;
		}

		.settings-form {
			grid-template-columns: 1fr;
		}

		.row-actions {
			display: grid;
			justify-content: stretch;
		}
	}
</style>
