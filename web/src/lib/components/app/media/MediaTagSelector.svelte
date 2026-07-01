<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import type { Tag } from '$lib/settings/types';

	interface Props {
		tags: Tag[];
		selectedTags: string[];
	}

	let { tags, selectedTags = $bindable() }: Props = $props();
	let tagInput = $state('');

	function toggleTag(name: string) {
		selectedTags = selectedTags.some((tag) => tag.toLowerCase() === name.toLowerCase())
			? selectedTags.filter((tag) => tag.toLowerCase() !== name.toLowerCase())
			: [...selectedTags, name];
	}

	function removeTag(name: string) {
		selectedTags = selectedTags.filter((tag) => tag.toLowerCase() !== name.toLowerCase());
	}

	function commitTagInput() {
		const name = tagInput.trim().replace(/\s+/g, ' ');
		if (!name || selectedTags.some((tag) => tag.toLowerCase() === name.toLowerCase())) {
			tagInput = '';
			return;
		}
		selectedTags = [...selectedTags, name];
		tagInput = '';
	}

	function handleTagKeydown(event: globalThis.KeyboardEvent) {
		if (event.key !== 'Enter' && event.key !== ',') {
			return;
		}
		event.preventDefault();
		commitTagInput();
	}
</script>

<div class="grid gap-2">
	<div class="grid gap-2">
		<Label>Tags</Label>
		<div class="flex min-h-10 flex-wrap items-center gap-2 rounded-md bg-muted/30 p-2">
			{#each selectedTags as tag (tag.toLowerCase())}
				<Button type="button" variant="secondary" size="xs" onclick={() => removeTag(tag)}>
					{tag}
				</Button>
			{/each}
			<Input
				bind:value={tagInput}
				class="h-7 min-w-32 flex-1 border-0 bg-transparent p-0 shadow-none focus-visible:ring-0"
				type="text"
				list="media-action-tag-options"
				maxlength={80}
				placeholder={selectedTags.length === 0 ? 'Add tag' : ''}
				autocomplete="off"
				onkeydown={handleTagKeydown}
				onblur={commitTagInput}
			/>
		</div>
		<datalist id="media-action-tag-options">
			{#each tags as tag (tag.id)}
				{#if !selectedTags.some((selected) => selected.toLowerCase() === tag.name.toLowerCase())}
					<option value={tag.name}></option>
				{/if}
			{/each}
		</datalist>
	</div>
	{#if tags.length > 0}
		<div class="flex flex-wrap gap-2" aria-label="Existing tags">
			{#each tags as tag (tag.id)}
				<Button
					type="button"
					variant={selectedTags.some(
						(selected) => selected.toLowerCase() === tag.name.toLowerCase()
					)
						? 'default'
						: 'outline'}
					size="xs"
					onclick={() => toggleTag(tag.name)}
				>
					{tag.name}
				</Button>
			{/each}
		</div>
	{/if}
</div>
