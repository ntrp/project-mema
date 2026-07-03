<script lang="ts">
	import { Label } from '$lib/components/ui/label';
	import { Slider } from '$lib/components/ui/slider';

	interface Props {
		label: string;
		value: number;
		min: number;
		max: number;
		step?: number;
		unit?: string;
		onChange: (_value: number) => void;
	}

	let { label, value, min, max, step = 1, unit = '', onChange }: Props = $props();
	let draftValue = $state<number | undefined>();
	const localValue = $derived(draftValue ?? clamp(value));

	const percentage = $derived(((localValue - min) / (max - min)) * 100);
	const rangeStyle = $derived(`left: 0%; width: ${Math.max(0, Math.min(100, percentage))}%`);

	function clamp(next: number) {
		return Math.max(min, Math.min(Number.isFinite(next) ? next : min, max));
	}

	function update(next: number[]) {
		draftValue = clamp(next[0] ?? min);
	}

	function commit(next = localValue) {
		const committed = clamp(next);
		draftValue = undefined;
		onChange(committed);
	}
</script>

<div class="grid gap-2">
	<Label>{label}</Label>
	<Slider
		value={[localValue]}
		{min}
		{max}
		{step}
		ariaLabel={label}
		rangeClass="hidden"
		{rangeStyle}
		onValueChange={update}
		onValueCommit={(next) => commit(next[0])}
	/>
	<span class="text-xs font-bold text-muted-foreground">{localValue}{unit}</span>
</div>
