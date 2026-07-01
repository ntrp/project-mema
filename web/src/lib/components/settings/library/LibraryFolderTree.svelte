<script lang="ts">
	import type { TreeNode } from './libraryFolderTree';

	interface Props {
		visibleNodes: TreeNode[];
		selectedPath: string;
		onToggle: (_path: string) => void;
	}

	let { visibleNodes, selectedPath = $bindable(), onToggle }: Props = $props();
</script>

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
				onclick={() => onToggle(node.path)}
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
				ondblclick={() => onToggle(node.path)}
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
