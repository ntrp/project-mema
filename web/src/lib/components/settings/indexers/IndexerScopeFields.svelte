<script lang="ts">
	import { Checkbox } from '$lib/components/ui/checkbox';
	import { Label } from '$lib/components/ui/label';
	import type { IndexerForm, IndexerMediaType, Tag } from '$lib/settings/types';

	interface Props {
		form: IndexerForm;
		tags: Tag[];
	}

	let { form = $bindable(), tags }: Props = $props();

	const mediaTypes: { value: IndexerMediaType; label: string }[] = [
		{ value: 'movie', label: 'Movies' },
		{ value: 'serie', label: 'Series' },
		{ value: 'anime', label: 'Anime' },
		{ value: 'audio', label: 'Audio' },
		{ value: 'book', label: 'Books' }
	];

	function toggleMediaType(value: IndexerMediaType) {
		const current = form.mediaTypeScopes ?? [];
		if (current.includes(value)) {
			if (current.length === 1) {
				return;
			}
			form.mediaTypeScopes = current.filter((scope) => scope !== value);
			return;
		}
		form.mediaTypeScopes = [...current, value];
	}

	function toggleTag(name: string) {
		const current = form.tagScopes ?? [];
		if (current.includes(name)) {
			form.tagScopes = current.filter((tag) => tag !== name);
			return;
		}
		form.tagScopes = [...current, name];
	}
</script>

<div class="grid gap-4 rounded-md border border-border p-3">
	<div class="space-y-2">
		<Label>Media scopes</Label>
		<div class="grid gap-2 [grid-template-columns:repeat(auto-fit,minmax(120px,1fr))]">
			{#each mediaTypes as option (option.value)}
				<label
					class="grid grid-cols-[18px_minmax(0,1fr)] items-center gap-2 rounded-md bg-muted p-2"
				>
					<Checkbox
						checked={(form.mediaTypeScopes ?? []).includes(option.value)}
						onclick={() => toggleMediaType(option.value)}
					/>
					<span class="truncate text-sm">{option.label}</span>
				</label>
			{/each}
		</div>
	</div>
	<div class="space-y-2">
		<Label>Tag scopes</Label>
		<div class="grid gap-2 [grid-template-columns:repeat(auto-fit,minmax(120px,1fr))]">
			{#each tags as tag (tag.id)}
				<label
					class="grid grid-cols-[18px_minmax(0,1fr)] items-center gap-2 rounded-md bg-muted p-2"
				>
					<Checkbox
						checked={(form.tagScopes ?? []).includes(tag.name)}
						onclick={() => toggleTag(tag.name)}
					/>
					<span class="truncate text-sm">{tag.name}</span>
				</label>
			{:else}
				<p class="m-0 text-sm text-muted-foreground">No tags configured</p>
			{/each}
		</div>
	</div>
</div>
