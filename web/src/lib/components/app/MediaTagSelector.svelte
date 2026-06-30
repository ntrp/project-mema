<script lang="ts">
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

<div class="tag-selector">
	<div class="tag-selector-header">
		<span>Tags</span>
		<div class="tag-input-box">
			{#each selectedTags as tag (tag.toLowerCase())}
				<button type="button" onclick={() => removeTag(tag)}>{tag}</button>
			{/each}
			<input
				bind:value={tagInput}
				type="text"
				list="media-action-tag-options"
				maxlength="80"
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
		<div class="tag-options" aria-label="Existing tags">
			{#each tags as tag (tag.id)}
				<button
					type="button"
					class:active-tag={selectedTags.some(
						(selected) => selected.toLowerCase() === tag.name.toLowerCase()
					)}
					onclick={() => toggleTag(tag.name)}
				>
					{tag.name}
				</button>
			{/each}
		</div>
	{/if}
</div>
