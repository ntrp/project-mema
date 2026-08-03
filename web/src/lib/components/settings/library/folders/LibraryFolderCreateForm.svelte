<script lang="ts">
	import FolderPlusIcon from '@lucide/svelte/icons/folder-plus';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';

	interface Props {
		name: string;
		disabled: boolean;
		creating: boolean;
		onCreate: () => void;
	}

	let { name = $bindable(), disabled, creating, onCreate }: Props = $props();
</script>

<form
	class="flex flex-wrap items-end justify-end gap-2.5"
	onsubmit={(event) => {
		event.preventDefault();
		onCreate();
	}}
>
	<label class="grid flex-[1_1_260px] gap-1.5">
		<span class="text-sm font-bold text-muted-foreground">Create under selected folder</span>
		<Input bind:value={name} placeholder="New folder name" maxlength={255} disabled={creating} />
	</label>
	<Button type="submit" variant="outline" disabled={disabled || !name.trim() || creating}>
		<FolderPlusIcon aria-hidden="true" />
		{creating ? 'Creating' : 'Create folder'}
	</Button>
</form>
