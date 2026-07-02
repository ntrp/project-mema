<script lang="ts">
	import SettingsFormModal from '$lib/components/settings/shared/SettingsFormModal.svelte';
	import { Button } from '$lib/components/ui/button';
	import { Card } from '$lib/components/ui/card';
	import { createLibraryFolderOption, listLibraryFolderOptions } from '$lib/settings/api';
	import { onMount } from 'svelte';
	import LibraryFolderCreateForm from './LibraryFolderCreateForm.svelte';
	import LibraryFolderTree from './LibraryFolderTree.svelte';
	import {
		createNode,
		flattenTree,
		type LibraryFolderPickerProps,
		type TreeNode
	} from './libraryFolderTree';

	let { initialPath, onClose, onUse }: LibraryFolderPickerProps = $props();
	let pickerLoading = $state(false);
	let pickerError = $state('');
	let rootPath = $state('');
	let rootParentPath = $state<string | undefined>();
	let selectedPath = $state('');
	let treeNodes = $state<Record<string, TreeNode>>({});
	let visibleNodes = $derived(flattenTree(treeNodes, rootPath));
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
</script>

<SettingsFormModal title="Select library folder" modalClass="grid gap-4" {onClose}>
	<div class="grid gap-4">
		<Card class="grid gap-2 p-3">
			<span class="text-xs font-black text-muted-foreground uppercase">Selected</span>
			<strong class="break-all text-sm text-foreground"
				>{selectedPath || 'No folder selected'}</strong
			>
		</Card>

		{#if pickerError}
			<p class="m-0 text-sm leading-6 text-muted-foreground">{pickerError}</p>
		{:else if pickerLoading}
			<p class="m-0 text-sm leading-6 text-muted-foreground">Loading folders</p>
		{:else}
			<div class="flex flex-wrap gap-2">
				<Button type="button" disabled={!selectedPath} onclick={useSelectedFolder}>
					Use selected folder
				</Button>
				{#if rootParentPath}
					<Button type="button" variant="outline" onclick={() => loadRoot(rootParentPath)}>
						Show parent
					</Button>
				{/if}
			</div>

			<LibraryFolderCreateForm
				bind:name={newFolderName}
				disabled={!selectedPath}
				creating={creatingFolder}
				onCreate={() => void createDirectory()}
			/>
			{#if createFolderError}
				<p class="m-0 text-sm font-semibold text-destructive">{createFolderError}</p>
			{/if}

			<LibraryFolderTree
				{visibleNodes}
				bind:selectedPath
				onToggle={(path) => void toggleNode(path)}
			/>
		{/if}
	</div>
</SettingsFormModal>
