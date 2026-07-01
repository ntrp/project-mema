<script lang="ts">
	interface Props {
		name: string;
		disabled: boolean;
		creating: boolean;
		onCreate: () => void;
	}

	let { name = $bindable(), disabled, creating, onCreate }: Props = $props();
</script>

<form
	class="folder-create-row"
	onsubmit={(event) => {
		event.preventDefault();
		onCreate();
	}}
>
	<label>
		<span>Create under selected folder</span>
		<input bind:value={name} placeholder="New folder name" maxlength="255" disabled={creating} />
	</label>
	<button type="submit" class="secondary" disabled={disabled || !name.trim() || creating}>
		{creating ? 'Creating' : 'Create folder'}
	</button>
</form>
