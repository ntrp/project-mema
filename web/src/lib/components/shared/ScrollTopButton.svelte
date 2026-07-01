<script lang="ts">
	import ChevronUpIcon from '@lucide/svelte/icons/chevron-up';
	import { onMount } from 'svelte';
	import { Button } from '$lib/components/ui/button';
	import * as Tooltip from '$lib/components/ui/tooltip';

	let visible = $state(false);

	onMount(() => {
		const handleScroll = () => {
			visible = window.scrollY > 700;
		};

		window.addEventListener('scroll', handleScroll, { passive: true });
		handleScroll();

		return () => window.removeEventListener('scroll', handleScroll);
	});

	function scrollToTop() {
		window.scrollTo({ top: 0, behavior: 'smooth' });
	}
</script>

{#if visible}
	<Tooltip.Root>
		<Tooltip.Trigger>
			{#snippet child({ props })}
				<Button
					{...props}
					type="button"
					variant="outline"
					size="icon"
					class="fixed right-6 bottom-6 z-[70] size-11 min-h-11 border-border bg-card/95 p-0 text-muted-foreground shadow-xl hover:border-primary/50 hover:bg-muted hover:text-primary"
					aria-label="Scroll to top"
					onclick={scrollToTop}
				>
					<ChevronUpIcon aria-hidden="true" />
				</Button>
			{/snippet}
		</Tooltip.Trigger>
		<Tooltip.Content>Scroll to top</Tooltip.Content>
	</Tooltip.Root>
{/if}
