<script lang="ts">
	import type { Snippet } from 'svelte';
	import { Button, type ButtonSize, type ButtonVariant } from '$lib/components/ui/button';
	import * as Dialog from '$lib/components/ui/dialog';
	import * as Tooltip from '$lib/components/ui/tooltip';

	interface Props {
		label: string;
		title: string;
		description: string;
		confirmLabel?: string;
		confirmingLabel?: string;
		cancelLabel?: string;
		variant?: ButtonVariant;
		size?: ButtonSize;
		disabled?: boolean;
		class?: string;
		tooltip?: string;
		onConfirm: () => void | Promise<void>;
		children?: Snippet;
	}

	let {
		label,
		title,
		description,
		confirmLabel = 'Delete',
		confirmingLabel = 'Deleting',
		cancelLabel = 'Cancel',
		variant = 'destructive',
		size = 'default',
		disabled = false,
		class: className,
		tooltip,
		onConfirm,
		children
	}: Props = $props();

	let open = $state(false);
	let confirming = $state(false);

	function requestConfirmation() {
		open = true;
	}

	async function confirmAction() {
		confirming = true;
		try {
			await onConfirm();
			open = false;
		} finally {
			confirming = false;
		}
	}
</script>

{#if tooltip}
	<Tooltip.Root>
		<Tooltip.Trigger>
			{#snippet child({ props })}
				<Button
					{...props}
					type="button"
					{variant}
					{size}
					class={className}
					aria-label={label}
					disabled={disabled || confirming}
					onclick={requestConfirmation}
				>
					{@render children?.()}
				</Button>
			{/snippet}
		</Tooltip.Trigger>
		<Tooltip.Content>{tooltip}</Tooltip.Content>
	</Tooltip.Root>
{:else}
	<Button
		type="button"
		{variant}
		{size}
		class={className}
		aria-label={label}
		disabled={disabled || confirming}
		onclick={requestConfirmation}
	>
		{@render children?.()}
	</Button>
{/if}

<Dialog.Root bind:open>
	<Dialog.Content class="w-[min(460px,calc(100vw-32px))]">
		<Dialog.Header>
			<Dialog.Title>{title}</Dialog.Title>
			<Dialog.Description>{description}</Dialog.Description>
		</Dialog.Header>
		<Dialog.Footer>
			<Button type="button" variant="outline" disabled={confirming} onclick={() => (open = false)}>
				{cancelLabel}
			</Button>
			<Button type="button" variant="destructive" disabled={confirming} onclick={confirmAction}>
				{confirming ? confirmingLabel : confirmLabel}
			</Button>
		</Dialog.Footer>
	</Dialog.Content>
</Dialog.Root>
