<script lang="ts">
	import PlusIcon from '@lucide/svelte/icons/plus';
	import TrashIcon from '@lucide/svelte/icons/trash-2';
	import { Button } from '$lib/components/ui/button';
	import { Checkbox } from '$lib/components/ui/checkbox';
	import * as Card from '$lib/components/ui/card';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import * as Select from '$lib/components/ui/select';
	import type {
		MediaProfileComponentTarget,
		MediaProfileComponentType,
		MediaProfileForm
	} from '$lib/settings/types';

	interface Props {
		form: MediaProfileForm;
		onChange: (_form: MediaProfileForm) => void;
	}

	let { form, onChange }: Props = $props();
	const targets = $derived(form.componentTargets ?? []);

	function patch(componentTargets: MediaProfileComponentTarget[]) {
		onChange({ ...form, componentTargets });
	}

	function add(componentType: MediaProfileComponentType) {
		patch([
			...targets,
			{
				componentType,
				required: true,
				source: componentType === 'subtitle' ? 'subtitleProvider' : 'release',
				fallbackBehavior: 'strict'
			}
		]);
	}

	function update(index: number, value: Partial<MediaProfileComponentTarget>) {
		patch(targets.map((target, row) => (row === index ? { ...target, ...value } : target)));
	}

	function remove(index: number) {
		patch(targets.filter((_, row) => row !== index));
	}

	function componentLabel(value: string) {
		if (value === 'audio') return 'Audio';
		if (value === 'subtitle') return 'Subtitle';
		return 'Video';
	}

	function sourceLabel(value: string) {
		if (value === 'subtitleProvider') return 'Subtitle provider';
		if (value === 'existing') return 'Existing file';
		return 'Release';
	}

	function fallbackLabel(value: string) {
		if (value === 'preferExisting') return 'Prefer existing';
		if (value === 'allowMissing') return 'Allow missing';
		return 'Strict';
	}
</script>

<Card.Root>
	<Card.Header><Card.Title>Component targets</Card.Title></Card.Header>
	<Card.Content class="mt-2 grid gap-3">
		<div class="flex flex-wrap gap-2">
			<Button type="button" variant="outline" size="sm" onclick={() => add('video')}>
				<PlusIcon aria-hidden="true" /> Video
			</Button>
			<Button type="button" variant="outline" size="sm" onclick={() => add('audio')}>
				<PlusIcon aria-hidden="true" /> Audio
			</Button>
			<Button type="button" variant="outline" size="sm" onclick={() => add('subtitle')}>
				<PlusIcon aria-hidden="true" /> Subtitle
			</Button>
		</div>

		{#each targets as target, index (target.id ?? `${target.componentType}-${index}`)}
			<div
				class="grid gap-2 rounded-md bg-muted/30 p-3 text-sm xl:grid-cols-[120px_110px_1fr_1fr_1fr_150px_140px_auto] xl:items-end"
			>
				<div class="grid gap-1">
					<Label>Type</Label>
					<Select.Root
						type="single"
						value={target.componentType}
						onValueChange={(value: string) =>
							update(index, { componentType: value as MediaProfileComponentType })}
					>
						<Select.Trigger>{componentLabel(target.componentType)}</Select.Trigger>
						<Select.Content>
							<Select.Item value="video" label="Video" />
							<Select.Item value="audio" label="Audio" />
							<Select.Item value="subtitle" label="Subtitle" />
						</Select.Content>
					</Select.Root>
				</div>
				<Label class="flex items-center gap-2 pb-2">
					<Checkbox
						checked={target.required}
						onCheckedChange={(checked) => update(index, { required: checked === true })}
					/>
					<span>Required</span>
				</Label>
				<div class="grid gap-1">
					<Label>Language</Label>
					<Input
						value={target.languageId ?? ''}
						disabled={target.componentType === 'video'}
						oninput={(event) => update(index, { languageId: event.currentTarget.value })}
					/>
				</div>
				<div class="grid gap-1">
					<Label>Codec</Label>
					<Input
						value={target.codec ?? ''}
						oninput={(event) => update(index, { codec: event.currentTarget.value })}
					/>
				</div>
				<div class="grid gap-1">
					<Label>Channels</Label>
					<Input
						value={target.channels ?? ''}
						disabled={target.componentType !== 'audio'}
						oninput={(event) => update(index, { channels: event.currentTarget.value })}
					/>
				</div>
				<div class="grid gap-1">
					<Label>Source</Label>
					<Select.Root
						type="single"
						value={target.source}
						onValueChange={(value: string) =>
							update(index, {
								source: value as MediaProfileComponentTarget['source']
							})}
					>
						<Select.Trigger>{sourceLabel(target.source)}</Select.Trigger>
						<Select.Content>
							<Select.Item value="release" label="Release" />
							<Select.Item value="subtitleProvider" label="Subtitle provider" />
							<Select.Item value="existing" label="Existing file" />
						</Select.Content>
					</Select.Root>
				</div>
				<div class="grid gap-1">
					<Label>Fallback</Label>
					<Select.Root
						type="single"
						value={target.fallbackBehavior}
						onValueChange={(value: string) =>
							update(index, {
								fallbackBehavior: value as MediaProfileComponentTarget['fallbackBehavior']
							})}
					>
						<Select.Trigger>{fallbackLabel(target.fallbackBehavior)}</Select.Trigger>
						<Select.Content>
							<Select.Item value="strict" label="Strict" />
							<Select.Item value="preferExisting" label="Prefer existing" />
							<Select.Item value="allowMissing" label="Allow missing" />
						</Select.Content>
					</Select.Root>
				</div>
				<Button type="button" variant="destructive" size="icon-sm" onclick={() => remove(index)}>
					<TrashIcon aria-label="Remove component target" />
				</Button>
			</div>
		{:else}
			<p class="m-0 text-sm text-muted-foreground">No component targets configured.</p>
		{/each}
	</Card.Content>
</Card.Root>
