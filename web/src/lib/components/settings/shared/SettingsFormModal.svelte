<script lang="ts">
	import XIcon from '@lucide/svelte/icons/x';
	import type { Snippet } from 'svelte';
	import { Button } from '$lib/components/ui/button';
	import * as Dialog from '$lib/components/ui/dialog';
	import { cn } from '$lib/utils';

	interface Props {
		title: string;
		onClose: () => void;
		children: Snippet;
		modalClass?: string;
	}

	let { title, onClose, children, modalClass = '' }: Props = $props();
	let open = $state(true);

	function handleOpenChange(nextOpen: boolean) {
		open = nextOpen;
		if (!nextOpen) {
			onClose();
		}
	}

	function stopScrollBubble(event: Event) {
		event.stopPropagation();
	}
</script>

<Dialog.Root bind:open onOpenChange={handleOpenChange}>
	<Dialog.Content
		preventScroll
		showCloseButton={false}
		class={cn(
			'w-fit min-w-[min(420px,calc(100vw-32px))] max-w-[calc(100vw-32px)] gap-5 overflow-hidden overscroll-contain p-0 sm:max-w-[calc(100vw-32px)]',
			modalClass
		)}
		aria-labelledby="settings-form-modal-title"
		onwheel={stopScrollBubble}
		ontouchmove={stopScrollBubble}
	>
		<Dialog.Header class="border-b border-border px-6 py-4">
			<div class="flex items-center justify-between gap-4">
				<Dialog.Title id="settings-form-modal-title">{title}</Dialog.Title>
				<Dialog.Close>
					{#snippet child({ props })}
						<Button {...props} variant="secondary" size="icon-sm" aria-label="Close">
							<XIcon aria-hidden="true" />
						</Button>
					{/snippet}
				</Dialog.Close>
			</div>
		</Dialog.Header>
		<div class="overflow-auto overscroll-contain px-6 py-5">
			{@render children()}
		</div>
	</Dialog.Content>
</Dialog.Root>
