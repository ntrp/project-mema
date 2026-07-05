<script lang="ts">
	import { Checkbox } from '$lib/components/ui/checkbox';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import type { MediaProfileLanguageScore } from '$lib/settings/types';

	interface Props {
		id: string;
		label: string;
		score?: MediaProfileLanguageScore;
		onToggle: (_language: string) => void;
		onScoreChange: (_language: string, _score: number) => void;
		onRequiredChange: (_language: string, _required: boolean) => void;
	}

	let { id, label, score, onToggle, onScoreChange, onRequiredChange }: Props = $props();
	const selected = $derived(score !== undefined);
</script>

<div class="grid gap-2 rounded-md bg-muted/20 p-2 sm:grid-cols-[1fr_120px_80px] sm:items-center">
	<Label class="flex items-center gap-2 text-sm">
		<Checkbox checked={selected} onCheckedChange={() => onToggle(id)} />
		<span>{label}</span>
	</Label>
	<Label class="flex items-center gap-2 text-sm">
		<Checkbox
			checked={score?.required === true}
			disabled={!selected}
			onCheckedChange={(checked) => onRequiredChange(id, checked === true)}
		/>
		<span>Required</span>
	</Label>
	<Input
		class="w-20"
		type="number"
		aria-label={`${label} score`}
		value={score?.score ?? 0}
		disabled={!selected}
		inputmode="numeric"
		oninput={(event) => onScoreChange(id, event.currentTarget.valueAsNumber)}
	/>
</div>
