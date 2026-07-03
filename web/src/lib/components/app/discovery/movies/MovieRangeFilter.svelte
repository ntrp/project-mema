<script lang="ts">
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import { Slider } from '$lib/components/ui/slider';

	interface Props {
		label: string;
		value: [number, number];
		min: number;
		max: number;
		step?: number;
		unit?: string;
		onChange: (_value: [number, number]) => void;
	}

	let { label, value, min, max, step = 1, unit = '', onChange }: Props = $props();
	const id = $props.id();
	let draftValue = $state<[number, number] | undefined>();
	const localValue = $derived(draftValue ?? normalise(value));

	function normalise(next: number[]) {
		const low = Math.max(min, Math.min(next[0] ?? min, max));
		const high = Math.max(low, Math.min(next[1] ?? max, max));
		return [low, high] as [number, number];
	}

	function updateLocal(next: number[]) {
		draftValue = normalise(next);
	}

	function commit(next = localValue) {
		const committed = normalise(next);
		draftValue = undefined;
		onChange(committed);
	}

	function updateField(index: 0 | 1, raw: string) {
		const next = [...localValue] as [number, number];
		const parsed = Number(raw);
		next[index] = Number.isFinite(parsed) ? parsed : index === 0 ? min : max;
		updateLocal(next);
	}
</script>

<div class="grid gap-2">
	<Label for={`${id}-min`}>{label}</Label>
	<Slider
		value={[localValue[0], localValue[1]]}
		{min}
		{max}
		{step}
		ariaLabel={label}
		onValueChange={updateLocal}
		onValueCommit={(next) => commit(normalise(next))}
	/>
	<div class="grid grid-cols-2 gap-2">
		<Input
			id={`${id}-min`}
			type="number"
			{min}
			{max}
			{step}
			value={localValue[0]}
			oninput={(event) => updateField(0, event.currentTarget.value)}
			onchange={() => commit()}
		/>
		<Input
			type="number"
			{min}
			{max}
			{step}
			value={localValue[1]}
			oninput={(event) => updateField(1, event.currentTarget.value)}
			onchange={() => commit()}
		/>
	</div>
	<span class="text-xs font-bold text-muted-foreground"
		>{localValue[0]}{unit} - {localValue[1]}{unit}</span
	>
</div>
