<script lang="ts">
	import TemplateTokenTextarea from '$lib/components/settings/library/TemplateTokenTextarea.svelte';
	import { Separator } from '$lib/components/ui/separator';
	import type { FileNamingSettingsRequest } from '$lib/settings/types';

	type TemplateField = keyof FileNamingSettingsRequest;

	interface Props {
		id: string;
		title: string;
		fields: { key: TemplateField; label: string }[];
		templates: FileNamingSettingsRequest;
		onChange: (_key: TemplateField, _value: string) => void;
		example: (_value: string) => string;
	}

	let { id, title, fields, templates, onChange, example }: Props = $props();
</script>

<section class="grid gap-3" aria-labelledby={id}>
	<div class="grid gap-2">
		<h3 {id} class="m-0 text-lg text-foreground">{title}</h3>
		<Separator />
	</div>
	<div class="grid gap-3">
		{#each fields as field (field.key)}
			<label class="grid gap-1.5">
				<span class="text-sm font-extrabold text-muted-foreground">{field.label}</span>
				<TemplateTokenTextarea
					value={templates[field.key]}
					onChange={(value) => onChange(field.key, value)}
				/>
				<span class="min-w-0 truncate font-mono text-xs text-muted-foreground">
					{example(templates[field.key])}
				</span>
			</label>
		{/each}
	</div>
</section>
