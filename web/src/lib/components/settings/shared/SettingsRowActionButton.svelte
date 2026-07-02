<script lang="ts">
	import EditIcon from '@lucide/svelte/icons/pencil';
	import RefreshCwIcon from '@lucide/svelte/icons/refresh-cw';
	import TrashIcon from '@lucide/svelte/icons/trash-2';
	import { Button } from '$lib/components/ui/button';

	type RowActionIcon = 'delete' | 'edit' | 'sync';

	interface Props {
		label: string;
		icon: RowActionIcon;
		variant?: 'outline' | 'destructive';
		disabled?: boolean;
		href?: string;
		onclick?: () => void | Promise<void>;
	}

	let { label, icon, variant = 'outline', disabled = false, href, onclick }: Props = $props();
	const Icon = $derived(icon === 'delete' ? TrashIcon : icon === 'sync' ? RefreshCwIcon : EditIcon);
</script>

<Button type="button" {href} {variant} size="icon-sm" aria-label={label} {disabled} {onclick}>
	<Icon aria-hidden="true" />
</Button>
