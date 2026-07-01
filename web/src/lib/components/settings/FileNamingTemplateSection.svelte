<script lang="ts">
	import TemplateTokenTextarea from '$lib/components/settings/TemplateTokenTextarea.svelte';
	import type { FileNamingSettingsRequest } from '$lib/settings/types';

	type TemplateField = keyof FileNamingSettingsRequest;

	interface Props {
		id: string;
		title: string;
		fields: { key: TemplateField; label: string }[];
		templates: FileNamingSettingsRequest;
		onChange: (_key: TemplateField, _value: string) => void;
	}

	let { id, title, fields, templates, onChange }: Props = $props();
</script>

<section class="grid gap-3 rounded-md border border-border bg-card p-3.5" aria-labelledby={id}>
	<h3 {id} class="m-0 text-lg text-foreground">{title}</h3>
	{#each fields as field (field.key)}
		<label class="grid gap-1.5">
			<span class="text-sm font-extrabold text-muted-foreground">{field.label}</span>
			<TemplateTokenTextarea
				value={templates[field.key]}
				onChange={(value) => onChange(field.key, value)}
			/>
		</label>
	{/each}
</section>
