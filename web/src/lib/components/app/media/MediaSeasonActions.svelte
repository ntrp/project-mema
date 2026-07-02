<script lang="ts">
	import SearchIcon from '@lucide/svelte/icons/search';
	import UserIcon from '@lucide/svelte/icons/user';
	import { Button } from '$lib/components/ui/button';
	import * as Tooltip from '$lib/components/ui/tooltip';

	interface Props {
		canManage: boolean;
		busy: boolean;
		onAutoSearch: () => void;
		onManualSearch: () => void;
	}

	let { canManage, busy, onAutoSearch, onManualSearch }: Props = $props();

	function stopActionClick(event: Event) {
		event.stopPropagation();
	}

	function stopActionKeydown(event: KeyboardEvent) {
		event.stopPropagation();
	}
</script>

<div
	role="presentation"
	class="flex shrink-0 items-center justify-end gap-2"
	onclick={stopActionClick}
	onkeydown={stopActionKeydown}
>
	<Tooltip.Root>
		<Tooltip.Trigger>
			{#snippet child({ props })}
				<Button
					{...props}
					type="button"
					variant="outline"
					size="icon-sm"
					aria-label="Automatic season search"
					disabled={!canManage || busy}
					onclick={onAutoSearch}
				>
					<SearchIcon aria-hidden="true" />
				</Button>
			{/snippet}
		</Tooltip.Trigger>
		<Tooltip.Content>Automatic season search</Tooltip.Content>
	</Tooltip.Root>
	<Tooltip.Root>
		<Tooltip.Trigger>
			{#snippet child({ props })}
				<Button
					{...props}
					type="button"
					variant="outline"
					size="icon-sm"
					aria-label="Manual season search"
					disabled={busy}
					onclick={onManualSearch}
				>
					<UserIcon aria-hidden="true" />
				</Button>
			{/snippet}
		</Tooltip.Trigger>
		<Tooltip.Content>Manual season search</Tooltip.Content>
	</Tooltip.Root>
</div>
