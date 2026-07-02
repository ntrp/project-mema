<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import { Textarea } from '$lib/components/ui/textarea';
	import type { LanguageForm as LanguageFormValue } from '$lib/settings/types';

	interface Props {
		form: LanguageFormValue;
		saving: boolean;
		onSave: (_event: SubmitEvent) => void | Promise<void>;
		onCancel: () => void;
	}

	let { form = $bindable(), saving, onSave, onCancel }: Props = $props();
</script>

<form class="grid gap-4" onsubmit={onSave}>
	<div class="grid gap-2">
		<Label for="language-code">ISO code</Label>
		<Input
			id="language-code"
			bind:value={form.code}
			type="text"
			maxlength={8}
			disabled={Boolean(form.originalCode)}
			required
		/>
	</div>
	<div class="grid gap-2">
		<Label for="language-display-name">Display name</Label>
		<Input
			id="language-display-name"
			bind:value={form.displayName}
			type="text"
			maxlength={120}
			required
		/>
	</div>
	<div class="grid gap-2">
		<Label for="language-aliases">Aliases</Label>
		<Textarea
			id="language-aliases"
			bind:value={form.aliasesText}
			rows={5}
			placeholder="DE, DEU, GER, DEUTSCH"
		/>
	</div>
	<div class="flex justify-end gap-2">
		<Button type="button" variant="outline" onclick={onCancel}>Cancel</Button>
		<Button type="submit" disabled={saving}>
			{saving ? 'Saving' : form.originalCode ? 'Update language' : 'Create language'}
		</Button>
	</div>
</form>
