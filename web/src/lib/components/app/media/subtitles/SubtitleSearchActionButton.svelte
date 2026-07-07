<script lang="ts">
	import DownloadIcon from '@lucide/svelte/icons/download';
	import RefreshCwIcon from '@lucide/svelte/icons/refresh-cw';
	import SearchIcon from '@lucide/svelte/icons/search';
	import { Button } from '$lib/components/ui/button';
	import * as Tooltip from '$lib/components/ui/tooltip';

	type Icon = 'download' | 'refresh' | 'search';

	interface Props {
		icon: Icon;
		label: string;
		disabled?: boolean;
		onClick: () => void;
	}

	let { icon, label, disabled = false, onClick }: Props = $props();
</script>

<Tooltip.Root>
	<Tooltip.Trigger>
		{#snippet child({ props })}
			<Button
				{...props}
				type="button"
				variant="outline"
				size="icon-sm"
				aria-label={label}
				{disabled}
				onclick={onClick}
			>
				{#if icon === 'refresh'}
					<RefreshCwIcon aria-hidden="true" />
				{:else if icon === 'search'}
					<SearchIcon aria-hidden="true" />
				{:else}
					<DownloadIcon aria-hidden="true" />
				{/if}
			</Button>
		{/snippet}
	</Tooltip.Trigger>
	<Tooltip.Content>{label}</Tooltip.Content>
</Tooltip.Root>
