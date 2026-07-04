<script lang="ts">
	import SettingsSelect from '$lib/components/settings/shared/SettingsSelect.svelte';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import { Switch } from '$lib/components/ui/switch';
	import type { IndexerCatalogEntry } from '$lib/settings/types';

	type IndexerField = IndexerCatalogEntry['fields'][number];

	interface Props {
		field: IndexerField;
		value: unknown;
		onValueChange: (_value: unknown) => void;
	}

	let { field, value, onValueChange }: Props = $props();
	const fieldId = $derived(`indexer-field-${field.name}`);
	const textValue = $derived(value == null ? '' : String(value));
	const checkedValue = $derived(value === true || value === 'true');
	const options = $derived(
		(field.selectOptions ?? []).map((option) => ({ value: option.value, label: option.name }))
	);

	function updateNumber(raw: string) {
		if (raw.trim() === '') {
			onValueChange('');
			return;
		}
		const parsed = field.isFloat ? Number.parseFloat(raw) : Number.parseInt(raw, 10);
		onValueChange(Number.isFinite(parsed) ? parsed : raw);
	}
</script>

{#if field.type === 'info'}
	<div class="rounded-md border border-border bg-muted/40 p-3 text-sm text-muted-foreground">
		<div class="font-bold text-foreground">{field.label}</div>
		{#if field.helpText}
			<p class="m-0 mt-1">{field.helpText}</p>
		{/if}
	</div>
{:else}
	<div class="space-y-2">
		<Label for={fieldId}>{field.label}{field.unit ? ` (${field.unit})` : ''}</Label>
		{#if field.type === 'checkbox'}
			<div class="flex items-center gap-3 py-2">
				<Switch id={fieldId} checked={checkedValue} onCheckedChange={onValueChange} />
				<span class="text-sm text-muted-foreground">{checkedValue ? 'Enabled' : 'Disabled'}</span>
			</div>
		{:else if field.type === 'select'}
			<SettingsSelect value={textValue} {options} {onValueChange} />
		{:else}
			<Input
				id={fieldId}
				value={textValue}
				type={field.type === 'password'
					? 'password'
					: field.type === 'number'
						? 'number'
						: field.type}
				step={field.type === 'number' && field.isFloat ? '0.01' : undefined}
				placeholder={field.placeholder}
				autocomplete={field.type === 'password' ? 'off' : undefined}
				oninput={(event) =>
					field.type === 'number'
						? updateNumber(event.currentTarget.value)
						: onValueChange(event.currentTarget.value)}
			/>
		{/if}
		{#if field.helpText}
			<p class="m-0 text-xs text-muted-foreground">{field.helpText}</p>
		{/if}
		{#if field.helpTextWarning}
			<p class="m-0 text-xs text-destructive">{field.helpTextWarning}</p>
		{/if}
		{#if field.helpLink}
			<a
				class="text-xs font-bold text-primary underline"
				href={field.helpLink}
				target="_blank"
				rel="noreferrer"
			>
				Help
			</a>
		{/if}
	</div>
{/if}
