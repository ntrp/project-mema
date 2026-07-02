<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';

	interface Props {
		pattern: string;
		clearing: boolean;
		onClearAll: () => void | Promise<void>;
		onClearPattern: (_pattern: string) => void | Promise<void>;
	}

	let { pattern = $bindable(), clearing, onClearAll, onClearPattern }: Props = $props();

	async function clearPattern(event: SubmitEvent) {
		event.preventDefault();
		const nextPattern = pattern.trim();
		if (!nextPattern) return;
		await onClearPattern(nextPattern);
		pattern = '';
	}
</script>

<form class="grid items-end gap-3 md:grid-cols-[minmax(0,1fr)_auto]" onsubmit={clearPattern}>
	<div class="grid gap-1.5">
		<Label>Reset by regex</Label>
		<Input bind:value={pattern} placeholder="discover:|details:123|matrix" autocomplete="off" />
	</div>
	<div class="flex flex-wrap justify-end gap-2">
		<Button type="submit" variant="destructive" disabled={clearing || pattern.trim().length === 0}>
			{clearing ? 'Resetting' : 'Reset matching'}
		</Button>
		<Button
			type="button"
			variant="destructive"
			disabled={clearing}
			onclick={() => void onClearAll()}
		>
			Reset all
		</Button>
	</div>
</form>
