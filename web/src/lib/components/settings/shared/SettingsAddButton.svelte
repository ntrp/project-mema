<script lang="ts">
	import PlusIcon from '@lucide/svelte/icons/plus';
	import { Button } from '$lib/components/ui/button';

	interface Props {
		label: string;
		onclick?: () => void;
		href?: string;
	}

	let { label, onclick, href }: Props = $props();
</script>

<Button type="button" {href} {onclick}>
	<PlusIcon aria-hidden="true" />
	<span>{label}</span>
</Button>
