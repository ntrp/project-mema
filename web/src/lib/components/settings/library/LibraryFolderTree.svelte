<script lang="ts">
	import ChevronDownIcon from '@lucide/svelte/icons/chevron-down';
	import ChevronRightIcon from '@lucide/svelte/icons/chevron-right';
	import { Button } from '$lib/components/ui/button';
	import { cn } from '$lib/utils';
	import type { TreeNode } from './libraryFolderTree';

	interface Props {
		visibleNodes: TreeNode[];
		selectedPath: string;
		onToggle: (_path: string) => void;
	}

	let { visibleNodes, selectedPath = $bindable(), onToggle }: Props = $props();

	const itemDepthClasses = [
		'pl-1.5',
		'pl-6',
		'pl-10',
		'pl-14',
		'pl-18',
		'pl-22',
		'pl-26',
		'pl-30',
		'pl-34'
	];
	const errorDepthClasses = [
		'ml-12',
		'ml-16',
		'ml-20',
		'ml-24',
		'ml-28',
		'ml-32',
		'ml-36',
		'ml-40',
		'ml-44'
	];

	function depthClass(classes: string[], depth: number) {
		return classes[Math.min(Math.max(depth, 0), classes.length - 1)];
	}
</script>

<div
	class="grid max-h-[420px] gap-1 overflow-auto rounded-md border p-2"
	role="tree"
	aria-label="Server folders"
>
	{#each visibleNodes as node (node.path)}
		<div
			class={cn(
				'grid grid-cols-[32px_minmax(0,1fr)] items-start gap-2 rounded-md pr-1.5 py-1',
				depthClass(itemDepthClasses, node.depth),
				selectedPath === node.path && 'bg-muted'
			)}
			role="treeitem"
			aria-selected={selectedPath === node.path}
			aria-expanded={node.children.length > 0 || !node.loaded ? node.expanded : undefined}
		>
			<Button
				type="button"
				variant="ghost"
				size="icon-sm"
				aria-label={node.expanded ? 'Collapse folder' : 'Expand folder'}
				disabled={node.loading}
				onclick={() => onToggle(node.path)}
			>
				{#if node.loading}
					...
				{:else if node.expanded}
					<ChevronDownIcon aria-hidden="true" />
				{:else}
					<ChevronRightIcon aria-hidden="true" />
				{/if}
			</Button>
			<Button
				type="button"
				variant="ghost"
				class="grid h-auto justify-items-start gap-1 px-2 py-1.5 text-left"
				ondblclick={() => onToggle(node.path)}
				onclick={() => (selectedPath = node.path)}
			>
				<span class="font-semibold">{node.name}</span>
				<small class="break-all text-muted-foreground">{node.path}</small>
			</Button>
		</div>
		{#if node.error}
			<p
				class={cn(
					'my-1 text-sm font-semibold text-destructive',
					depthClass(errorDepthClasses, node.depth)
				)}
			>
				{node.error}
			</p>
		{/if}
	{/each}
</div>
