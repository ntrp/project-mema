import { createMutation, useQueryClient } from '@tanstack/svelte-query';
import { createLibraryFolderOption, listLibraryFolderOptions } from '$lib/settings/api';
import type { LibraryFolderOptionListResponse } from '$lib/settings/types';
import { createNode, flattenTree, folderPathChain, type TreeNode } from './libraryFolderTree';

export const folderBrowserKeys = {
	all: ['settings', 'folder-browser'] as const,
	path: (path?: string) => [...folderBrowserKeys.all, path ?? ''] as const
};

export function createFolderBrowser(initialPath?: string) {
	const resources = createFolderBrowserResources();
	let loading = $state(false);
	let error = $state('');
	let selectedPath = $state('');
	let treeNodes = $state.raw<Record<string, TreeNode>>({});
	let newFolderName = $state('');
	let creatingFolder = $state(false);
	let createFolderError = $state('');
	let visibleNodes = $derived(flattenTree(treeNodes, '/'));

	async function initialize() {
		loading = true;
		error = '';
		try {
			const preferred = await loadPreferredPath();
			const targetPath = preferred.currentPath;
			const root = targetPath === '/' ? preferred : await resources.load('/');
			setRoot(root);
			for (const path of folderPathChain(targetPath).slice(1)) {
				if (!(await loadNode(path, path === targetPath ? preferred : undefined))) break;
			}
			selectedPath = targetPath;
		} catch (cause) {
			error = cause instanceof Error ? cause.message : 'Could not load folders';
		} finally {
			loading = false;
		}
	}

	async function loadPreferredPath() {
		try {
			return await resources.load(initialPath || undefined);
		} catch (cause) {
			if (!initialPath) throw cause;
			return resources.load();
		}
	}

	function setRoot(response: LibraryFolderOptionListResponse) {
		const root = createNode({ name: '/', path: '/' }, undefined, 0);
		root.expanded = true;
		root.loaded = true;
		root.children = response.entries.map((entry) => entry.path);
		const nodes: Record<string, TreeNode> = { '/': root };
		for (const entry of response.entries) nodes[entry.path] = createNode(entry, '/', 1);
		treeNodes = nodes;
	}

	async function loadNode(path: string, known?: LibraryFolderOptionListResponse) {
		const node = treeNodes[path];
		if (!node) return false;
		treeNodes = { ...treeNodes, [path]: { ...node, loading: true, error: '' } };
		try {
			const response = known ?? (await resources.load(path));
			const children = response.entries.map((entry) => entry.path);
			const nodes = {
				...treeNodes,
				[path]: { ...treeNodes[path], expanded: true, loaded: true, loading: false, children }
			};
			for (const entry of response.entries) {
				const existing = nodes[entry.path];
				nodes[entry.path] = existing ?? createNode(entry, path, node.depth + 1);
			}
			treeNodes = nodes;
			return true;
		} catch (cause) {
			const message = cause instanceof Error ? cause.message : 'Could not load folder';
			treeNodes = { ...treeNodes, [path]: { ...treeNodes[path], loading: false, error: message } };
			return false;
		}
	}

	async function toggleNode(path: string) {
		const node = treeNodes[path];
		if (!node) return;
		if (node.expanded || node.loaded) {
			treeNodes = { ...treeNodes, [path]: { ...node, expanded: !node.expanded } };
			return;
		}
		await loadNode(path);
	}

	async function createDirectory() {
		const name = newFolderName.trim();
		if (!selectedPath || !name) return;
		creatingFolder = true;
		createFolderError = '';
		try {
			const created = await resources.create.mutateAsync({ parentPath: selectedPath, name });
			newFolderName = '';
			await resources.invalidate(selectedPath);
			await loadNode(selectedPath);
			selectedPath = created.path;
		} catch (cause) {
			createFolderError = cause instanceof Error ? cause.message : 'Could not create folder';
		} finally {
			creatingFolder = false;
		}
	}

	return {
		get loading() {
			return loading;
		},
		get error() {
			return error;
		},
		get selectedPath() {
			return selectedPath;
		},
		set selectedPath(path: string) {
			selectedPath = path;
		},
		get visibleNodes() {
			return visibleNodes;
		},
		get newFolderName() {
			return newFolderName;
		},
		set newFolderName(name: string) {
			newFolderName = name;
		},
		get creatingFolder() {
			return creatingFolder;
		},
		get createFolderError() {
			return createFolderError;
		},
		initialize,
		toggleNode,
		createDirectory
	};
}

function createFolderBrowserResources() {
	const client = useQueryClient();
	return {
		load: (path?: string) =>
			client.fetchQuery({
				queryKey: folderBrowserKeys.path(path),
				queryFn: () => listLibraryFolderOptions(path),
				staleTime: 30_000
			}),
		create: createMutation(() => ({
			mutationFn: ({ parentPath, name }: { parentPath: string; name: string }) =>
				createLibraryFolderOption(parentPath, name)
		})),
		invalidate: (path: string) =>
			client.invalidateQueries({ queryKey: folderBrowserKeys.path(path) })
	};
}
