<script lang="ts">
	import BracesIcon from '@lucide/svelte/icons/braces';
	import CopyIcon from '@lucide/svelte/icons/copy';
	import CheckIcon from '@lucide/svelte/icons/check';
	import { Button } from '$lib/components/ui/button';
	import * as Dialog from '$lib/components/ui/dialog';

	interface Props {
		value: string;
		success: boolean;
	}

	let { value, success }: Props = $props();
	let copied = $state(false);
	let open = $state(false);
	const formatted = $derived(formatJson(value));

	async function copyResponse() {
		await globalThis.navigator.clipboard.writeText(formatted);
		copied = true;
		globalThis.setTimeout(() => {
			copied = false;
		}, 1200);
	}

	function formatJson(raw: string) {
		try {
			return JSON.stringify(JSON.parse(raw), null, 2);
		} catch {
			return raw;
		}
	}
</script>

<div class="flex items-center justify-end gap-1">
	<Button
		type="button"
		variant="ghost"
		size="icon-sm"
		aria-label="Show response JSON"
		class={success ? '' : 'text-destructive'}
		onclick={() => (open = true)}
	>
		<BracesIcon aria-hidden="true" />
	</Button>
</div>

<Dialog.Root bind:open>
	<Dialog.Content
		class="grid max-h-[min(720px,calc(100vh-32px))] w-[min(920px,calc(100vw-32px))] grid-rows-[auto_minmax(0,1fr)_auto] gap-4 sm:max-w-none"
	>
		<Dialog.Header>
			<Dialog.Title>Response JSON</Dialog.Title>
			<Dialog.Description>
				{success ? 'Provider response payload' : 'Provider error response'}
			</Dialog.Description>
		</Dialog.Header>

		<div class="min-h-0 overflow-auto rounded-md border border-border bg-muted/30 p-3">
			<pre class="m-0 whitespace-pre-wrap wrap-break-word text-xs leading-5">{formatted}</pre>
		</div>

		<Dialog.Footer>
			<Button type="button" variant="outline" onclick={() => void copyResponse()}>
				{#if copied}
					<CheckIcon aria-hidden="true" />
					Copied
				{:else}
					<CopyIcon aria-hidden="true" />
					Copy
				{/if}
			</Button>
		</Dialog.Footer>
	</Dialog.Content>
</Dialog.Root>
