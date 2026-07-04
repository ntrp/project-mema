<script lang="ts">
	import EditIcon from '@lucide/svelte/icons/pencil';
	import RefreshCwIcon from '@lucide/svelte/icons/refresh-cw';
	import TrashIcon from '@lucide/svelte/icons/trash-2';
	import ConfirmActionButton from '$lib/components/shared/ConfirmActionButton.svelte';
	import { Button } from '$lib/components/ui/button';

	type RowActionIcon = 'delete' | 'edit' | 'sync';

	interface Props {
		label: string;
		icon: RowActionIcon;
		variant?: 'outline' | 'destructive';
		disabled?: boolean;
		href?: string;
		confirmTitle?: string;
		confirmDescription?: string;
		confirmLabel?: string;
		confirmingLabel?: string;
		onclick?: () => void | Promise<void>;
	}

	let {
		label,
		icon,
		variant = 'outline',
		disabled = false,
		href,
		confirmTitle,
		confirmDescription,
		confirmLabel,
		confirmingLabel,
		onclick
	}: Props = $props();
	const Icon = $derived(icon === 'delete' ? TrashIcon : icon === 'sync' ? RefreshCwIcon : EditIcon);

	function runAction() {
		return onclick?.();
	}
</script>

{#if confirmTitle && confirmDescription}
	<ConfirmActionButton
		{label}
		title={confirmTitle}
		description={confirmDescription}
		{confirmLabel}
		{confirmingLabel}
		{variant}
		size="icon-sm"
		{disabled}
		onConfirm={runAction}
	>
		<Icon aria-hidden="true" />
	</ConfirmActionButton>
{:else}
	<Button type="button" {href} {variant} size="icon-sm" aria-label={label} {disabled} {onclick}>
		<Icon aria-hidden="true" />
	</Button>
{/if}
