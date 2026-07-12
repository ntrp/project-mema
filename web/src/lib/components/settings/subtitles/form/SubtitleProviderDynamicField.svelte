<script lang="ts">
	import SecretInput from '$lib/components/settings/shared/SecretInput.svelte';
	import { Checkbox } from '$lib/components/ui/checkbox';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import type { components } from '$lib/api/generated/schema';
	import type { SubtitleProviderForm } from '$lib/settings/formTypes';
	type SubtitleProviderField = components['schemas']['SubtitleProviderField'];
	type SettingValue = components['schemas']['SubtitleProviderSettingValue'];

	interface Props {
		field: SubtitleProviderField;
		form: SubtitleProviderForm;
		onChange: (_form: SubtitleProviderForm) => void;
	}

	let { field, form, onChange }: Props = $props();
	const current = $derived(form.settings?.[field.key] ?? {});
	const savedSecret = $derived(
		(field.secret && form.secretFieldsSet?.includes(field.key)) ||
			(field.key === 'apiKey' && form.apiKeySet) ||
			(field.key === 'password' && form.passwordSet)
	);

	function stringValue() {
		return current.stringValue ?? legacyTextValue(field.key) ?? '';
	}

	function legacyTextValue(key: string) {
		if (key === 'baseUrl') return form.baseUrl;
		if (key === 'username') return form.username;
		return undefined;
	}

	function updateSetting(value: SettingValue) {
		const settings = { ...(form.settings ?? {}), [field.key]: value };
		onChange({ ...form, ...legacyPatch(value), settings });
	}

	function updateSecret(value: string) {
		const secretSettings = { ...(form.secretSettings ?? {}), [field.key]: value };
		const clearSecretFields = (form.clearSecretFields ?? []).filter((item) => item !== field.key);
		if (value === '' && savedSecret) clearSecretFields.push(field.key);
		onChange({ ...form, ...legacySecretPatch(value), secretSettings, clearSecretFields });
	}

	function legacyPatch(value: SettingValue) {
		if (field.key === 'baseUrl') return { baseUrl: value.stringValue ?? '' };
		if (field.key === 'username') return { username: value.stringValue ?? '' };
		return {};
	}

	function legacySecretPatch(value: string) {
		if (field.key === 'apiKey') return { apiKey: value };
		if (field.key === 'password') return { password: value };
		return {};
	}
</script>

{#if field.type === 'switch'}
	<div class="flex min-h-9 items-center gap-2 self-end">
		<Checkbox
			id={`subtitle-field-${field.key}`}
			checked={current.booleanValue === true}
			onCheckedChange={(checked) => updateSetting({ booleanValue: checked === true })}
		/>
		<Label for={`subtitle-field-${field.key}`}>{field.label}</Label>
	</div>
{:else if field.type === 'select'}
	<label class="grid gap-1.5">
		<span class="text-sm font-bold text-muted-foreground">{field.label}</span>
		<select
			class="h-9 rounded-md border border-input bg-background px-3 text-sm"
			value={stringValue()}
			required={field.required}
			onchange={(event) => updateSetting({ stringValue: event.currentTarget.value })}
		>
			{#each field.options ?? [] as option (option)}
				<option value={option}>{option}</option>
			{/each}
		</select>
	</label>
{:else if field.type === 'chips'}
	<label class="grid gap-1.5 md:col-span-2">
		<span class="text-sm font-bold text-muted-foreground">{field.label}</span>
		<Input
			value={(current.stringValues ?? []).join(', ')}
			placeholder="Comma-separated values"
			oninput={(event) =>
				updateSetting({
					stringValues: event.currentTarget.value
						.split(',')
						.map((item) => item.trim())
						.filter(Boolean)
				})}
		/>
	</label>
{:else if field.type !== 'action'}
	<label class="grid gap-1.5">
		<span class="text-sm font-bold text-muted-foreground">{field.label}</span>
		{#if field.secret}
			<SecretInput
				value={form.secretSettings?.[field.key] ?? ''}
				placeholder={savedSecret ? 'Saved secret' : ''}
				autocomplete="off"
				onValueChange={updateSecret}
			/>
		{:else}
			<Input
				value={stringValue()}
				required={field.required}
				maxlength={2000}
				oninput={(event) => updateSetting({ stringValue: event.currentTarget.value })}
			/>
		{/if}
	</label>
{/if}
