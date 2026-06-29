<script lang="ts">
	import { createLibraryFolderOption, listLibraryFolderOptions } from '$lib/settings/api';
	import type { LibraryFolderOption } from '$lib/settings/types';
	import { onMount } from 'svelte';

	interface TreeNode {
		name: string;
		path: string;
		parentPath?: string;
		depth: number;
		expanded: boolean;
		loaded: boolean;
		loading: boolean;
		error: string;
		children: string[];
	}

	interface Props {
		initialPath?: string;
		onClose: () => void;
		onUse: (_path: string) => void;
	}

	let { initialPath, onClose, onUse }: Props = $props();
	let pickerLoading = $state(false);
	let pickerError = $state('');
	let rootPath = $state('');
	let rootParentPath = $state<string | undefined>();
	let selectedPath = $state('');
	let treeNodes = $state<Record<string, TreeNode>>({});
	let visibleNodes = $derived(flattenTree());
	let newFolderName = $state('');
	let creatingFolder = $state(false);
	let createFolderError = $state('');

	onMount(() => {
		void loadRoot(initialPath || undefined);
	});

	async function loadRoot(path?: string) {
		pickerLoading = true;
		pickerError = '';
		createFolderError = '';
		try {
			const response = await listLibraryFolderOptions(path);
			const root = createNode(
				{ name: response.currentPath, path: response.currentPath },
				undefined,
				0
			);
			root.expanded = true;
			root.loaded = true;
			root.children = response.entries.map((entry) => entry.path);

			const nextNodes: Record<string, TreeNode> = { [root.path]: root };
			for (const entry of response.entries) {
				nextNodes[entry.path] = createNode(entry, root.path, 1);
			}
			treeNodes = nextNodes;
			rootPath = response.currentPath;
			rootParentPath = response.parentPath;
			selectedPath = response.currentPath;
		} catch (error) {
			if (path) {
				await loadRoot();
				return;
			}
			pickerError = error instanceof Error ? error.message : 'Could not load folders';
		} finally {
			pickerLoading = false;
		}
	}

	async function refreshNode(path: string, nextSelectedPath?: string) {
		const node = treeNodes[path];
		if (!node) {
			await loadRoot(path);
			if (nextSelectedPath) {
				selectedPath = nextSelectedPath;
			}
			return;
		}

		treeNodes = { ...treeNodes, [path]: { ...node, loading: true, error: '' } };
		try {
			const response = await listLibraryFolderOptions(path);
			const children = response.entries.map((entry) => entry.path);
			const nextNodes = {
				...treeNodes,
				[path]: {
					...treeNodes[path],
					expanded: true,
					loaded: true,
					loading: false,
					error: '',
					children
				}
			};
			for (const entry of response.entries) {
				nextNodes[entry.path] = createNode(entry, path, node.depth + 1);
			}
			treeNodes = nextNodes;
			if (nextSelectedPath) {
				selectedPath = nextSelectedPath;
			}
		} catch (error) {
			const message = error instanceof Error ? error.message : 'Could not load folder';
			treeNodes = {
				...treeNodes,
				[path]: { ...treeNodes[path], loading: false, error: message }
			};
		}
	}

	async function toggleNode(path: string) {
		const node = treeNodes[path];
		if (!node) {
			return;
		}
		if (node.expanded || node.loaded) {
			treeNodes = { ...treeNodes, [path]: { ...node, expanded: !node.expanded } };
			return;
		}
		await refreshNode(path);
	}

	async function createDirectory() {
		const parentPath = selectedPath || rootPath;
		const name = newFolderName.trim();
		if (!parentPath || !name) {
			return;
		}

		creatingFolder = true;
		createFolderError = '';
		try {
			const created = await createLibraryFolderOption(parentPath, name);
			newFolderName = '';
			await refreshNode(parentPath, created.path);
		} catch (error) {
			createFolderError = error instanceof Error ? error.message : 'Could not create folder';
		} finally {
			creatingFolder = false;
		}
	}

	function useSelectedFolder() {
		if (selectedPath) {
			onUse(selectedPath);
		}
	}

	function createNode(
		entry: LibraryFolderOption,
		parentPath: string | undefined,
		depth: number
	): TreeNode {
		return {
			name: entry.name,
			path: entry.path,
			parentPath,
			depth,
			expanded: false,
			loaded: false,
			loading: false,
			error: '',
			children: []
		};
	}

	function flattenTree() {
		const root = treeNodes[rootPath];
		if (!root) {
			return [];
		}

		const flattened: TreeNode[] = [];
		appendVisible(root, flattened);
		return flattened;
	}

	function appendVisible(node: TreeNode, flattened: TreeNode[]) {
		flattened.push(node);
		if (!node.expanded) {
			return;
		}
		for (const childPath of node.children) {
			const child = treeNodes[childPath];
			if (child) {
				appendVisible(child, flattened);
			}
		}
	}

	function handleWindowKeydown(event: { key: string }) {
		if (event.key === 'Escape') {
			onClose();
		}
	}
</script>

<svelte:window onkeydown={handleWindowKeydown} />

<div class="modal-backdrop" role="presentation" onclick={onClose}>
	<div
		class="modal-shell folder-picker-modal"
		role="dialog"
		aria-modal="true"
		aria-labelledby="folder-picker-title"
		tabindex="-1"
		onclick={(event) => event.stopPropagation()}
		onkeydown={(event) => event.stopPropagation()}
	>
		<div class="section-heading">
			<div>
				<p class="section-kicker">Folder picker</p>
				<h3 id="folder-picker-title">Select library folder</h3>
			</div>
			<button type="button" class="secondary" onclick={onClose}>Close</button>
		</div>

		<div class="folder-picker-selected">
			<span>Selected</span>
			<strong>{selectedPath || 'No folder selected'}</strong>
		</div>

		{#if pickerError}
			<p class="empty">{pickerError}</p>
		{:else if pickerLoading}
			<p class="muted">Loading folders</p>
		{:else}
			<div class="folder-picker-toolbar">
				<button type="button" disabled={!selectedPath} onclick={useSelectedFolder}>
					Use selected folder
				</button>
				{#if rootParentPath}
					<button type="button" class="secondary" onclick={() => loadRoot(rootParentPath)}>
						Show parent
					</button>
				{/if}
			</div>

			<form
				class="folder-create-row"
				onsubmit={(event) => {
					event.preventDefault();
					void createDirectory();
				}}
			>
				<label>
					<span>Create under selected folder</span>
					<input
						bind:value={newFolderName}
						placeholder="New folder name"
						maxlength="255"
						disabled={creatingFolder}
					/>
				</label>
				<button
					type="submit"
					class="secondary"
					disabled={!selectedPath || !newFolderName.trim() || creatingFolder}
				>
					{creatingFolder ? 'Creating' : 'Create folder'}
				</button>
			</form>
			{#if createFolderError}
				<p class="folder-tree-error">{createFolderError}</p>
			{/if}

			<div class="folder-tree" role="tree" aria-label="Server folders">
				{#each visibleNodes as node (node.path)}
					<div
						class="folder-tree-row"
						class:selected={selectedPath === node.path}
						role="treeitem"
						aria-selected={selectedPath === node.path}
						aria-expanded={node.children.length > 0 || !node.loaded ? node.expanded : undefined}
						style={`--depth: ${node.depth}`}
					>
						<button
							type="button"
							class="folder-tree-toggle"
							aria-label={node.expanded ? 'Collapse folder' : 'Expand folder'}
							disabled={node.loading}
							onclick={() => toggleNode(node.path)}
						>
							{#if node.loading}
								...
							{:else if node.expanded}
								v
							{:else}
								&gt;
							{/if}
						</button>
						<button
							type="button"
							class="folder-tree-label"
							ondblclick={() => toggleNode(node.path)}
							onclick={() => (selectedPath = node.path)}
						>
							<span>{node.name}</span>
							<small>{node.path}</small>
						</button>
					</div>
					{#if node.error}
						<p class="folder-tree-error" style={`--depth: ${node.depth + 1}`}>{node.error}</p>
					{/if}
				{/each}
			</div>
		{/if}
	</div>
</div>
