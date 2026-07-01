import type { LibraryFolderOption } from '$lib/settings/types';

export interface TreeNode {
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

export interface LibraryFolderPickerProps {
	initialPath?: string;
	onClose: () => void;
	onUse: (_path: string) => void;
}

export function createNode(
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

export function flattenTree(treeNodes: Record<string, TreeNode>, rootPath: string) {
	const root = treeNodes[rootPath];
	if (!root) {
		return [];
	}

	const flattened: TreeNode[] = [];
	appendVisible(treeNodes, root, flattened);
	return flattened;
}

function appendVisible(treeNodes: Record<string, TreeNode>, node: TreeNode, flattened: TreeNode[]) {
	flattened.push(node);
	if (!node.expanded) {
		return;
	}
	for (const childPath of node.children) {
		const child = treeNodes[childPath];
		if (child) {
			appendVisible(treeNodes, child, flattened);
		}
	}
}
