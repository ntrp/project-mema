<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import * as Select from '$lib/components/ui/select';
	import type {
		ReleaseFilters,
		ReleaseSourceFilter
	} from '$lib/components/app/media/release-display/releaseSearchResults';

	interface Props {
		filters: ReleaseFilters;
		qualityOptions: string[];
		onChange: (_filters: ReleaseFilters) => void;
		onReset: () => void;
	}

	let { filters, qualityOptions, onChange, onReset }: Props = $props();

	function patch(patchValue: Partial<ReleaseFilters>) {
		onChange(normalizedFilters({ ...filters, ...patchValue }));
	}

	function protocolLabel(value: ReleaseSourceFilter) {
		if (value === 'all') return 'All';
		return value.toUpperCase();
	}

	function scrollNumber(
		event: globalThis.WheelEvent,
		key: 'minSize' | 'maxSize' | 'minScore' | 'maxScore',
		baseStep: number,
		shiftStep: number,
		minimum?: number
	) {
		event.preventDefault();
		const current = Number(filters[key]);
		const direction = event.deltaY < 0 ? 1 : -1;
		const step = event.shiftKey ? shiftStep : baseStep;
		let next = (Number.isFinite(current) ? current : 0) + direction * step;
		if (minimum !== undefined) {
			next = Math.max(minimum, next);
		}
		patch({ [key]: formattedNumber(next) });
	}

	function normalizedFilters(next: ReleaseFilters) {
		return normalizePair(
			normalizePair(next, 'minSize', 'maxSize', 0.1),
			'minScore',
			'maxScore',
			100
		);
	}

	function normalizePair(
		filtersValue: ReleaseFilters,
		minKey: 'minSize' | 'minScore',
		maxKey: 'maxSize' | 'maxScore',
		step: number
	) {
		const min = finiteNumber(filtersValue[minKey]);
		const max = finiteNumber(filtersValue[maxKey]);
		if (min === undefined || max === undefined || max > min) {
			return filtersValue;
		}
		return { ...filtersValue, [maxKey]: formattedNumber(min + step) };
	}

	function finiteNumber(value: string) {
		if (value.trim() === '') return undefined;
		const parsed = Number(value);
		return Number.isFinite(parsed) ? parsed : undefined;
	}

	function formattedNumber(value: number) {
		return String(Number(value.toFixed(2)));
	}
</script>

<div class="grid gap-3 rounded-md border border-border p-3">
	<div class="grid gap-3 sm:grid-cols-[160px_1fr_1fr_220px_auto] sm:items-end">
		<div class="grid gap-1.5">
			<Label>Protocol</Label>
			<Select.Root
				type="single"
				value={filters.source}
				onValueChange={(value) => patch({ source: value as ReleaseSourceFilter })}
			>
				<Select.Trigger class="w-full">{protocolLabel(filters.source)}</Select.Trigger>
				<Select.Content>
					<Select.Item value="all" label="All" />
					<Select.Item value="nzb" label="NZB" />
					<Select.Item value="torrent" label="TORRENT" />
				</Select.Content>
			</Select.Root>
		</div>
		<div class="grid gap-1.5">
			<Label>Size GiB</Label>
			<div class="grid grid-cols-2 gap-2">
				<Input
					type="number"
					min="0"
					step="0.1"
					placeholder="Min"
					value={filters.minSize}
					oninput={(event) => patch({ minSize: event.currentTarget.value })}
					onwheel={(event) => scrollNumber(event, 'minSize', 0.1, 1, 0)}
				/>
				<Input
					type="number"
					min="0"
					step="0.1"
					placeholder="Max"
					value={filters.maxSize}
					oninput={(event) => patch({ maxSize: event.currentTarget.value })}
					onwheel={(event) => scrollNumber(event, 'maxSize', 0.1, 1, 0)}
				/>
			</div>
		</div>
		<div class="grid gap-1.5">
			<Label>Score</Label>
			<div class="grid grid-cols-2 gap-2">
				<Input
					type="number"
					step="1"
					placeholder="Min"
					value={filters.minScore}
					oninput={(event) => patch({ minScore: event.currentTarget.value })}
					onwheel={(event) => scrollNumber(event, 'minScore', 100, 1000)}
				/>
				<Input
					type="number"
					step="1"
					placeholder="Max"
					value={filters.maxScore}
					oninput={(event) => patch({ maxScore: event.currentTarget.value })}
					onwheel={(event) => scrollNumber(event, 'maxScore', 100, 1000)}
				/>
			</div>
		</div>
		<div class="grid gap-1.5">
			<Label>Quality</Label>
			<Select.Root
				type="single"
				value={filters.quality}
				onValueChange={(value) => patch({ quality: value })}
			>
				<Select.Trigger class="w-full">
					{filters.quality === 'all' ? 'All qualities' : filters.quality}
				</Select.Trigger>
				<Select.Content>
					<Select.Item value="all" label="All qualities" />
					{#each qualityOptions as quality (quality)}
						<Select.Item value={quality} label={quality} />
					{/each}
				</Select.Content>
			</Select.Root>
		</div>
		<Button type="button" variant="outline" onclick={onReset}>Reset</Button>
	</div>
</div>
