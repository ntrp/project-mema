<script lang="ts">
	import PlusIcon from '@lucide/svelte/icons/plus';
	import { Button } from '$lib/components/ui/button';

	interface Props {
		label: string;
		onclick: () => void;
	}

	let { label, onclick }: Props = $props();
</script>

<Button type="button" {onclick}>
	<PlusIcon aria-hidden="true" />
	<span>{label}</span>
</Button>
