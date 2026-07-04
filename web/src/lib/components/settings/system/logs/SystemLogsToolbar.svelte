<script lang="ts">
	import ConfirmActionButton from '$lib/components/shared/ConfirmActionButton.svelte';
	import SettingsSelect from '$lib/components/settings/shared/SettingsSelect.svelte';
	import { Button } from '$lib/components/ui/button';
	import { Label } from '$lib/components/ui/label';
	import type { SystemLogLevel } from '$lib/settings/types';

	interface Props {
		level: SystemLogLevel;
		levelOptions: { value: SystemLogLevel; label: string }[];
		loading: boolean;
		saving: boolean;
		followLogs: boolean;
		onClearLogs: () => void;
		onEnableFollow: () => void;
		onLevelChange: (_value: string) => void | Promise<void>;
	}

	let {
		level,
		levelOptions,
		loading,
		saving,
		followLogs,
		onClearLogs,
		onEnableFollow,
		onLevelChange
	}: Props = $props();
</script>

<div class="flex flex-wrap items-end gap-3">
	<ConfirmActionButton
		label="Clear logs"
		title="Clear logs"
		description="Clear the visible log buffer?"
		confirmLabel="Clear logs"
		variant="outline"
		size="sm"
		onConfirm={onClearLogs}
	>
		Clear logs
	</ConfirmActionButton>
	<Button
		type="button"
		variant="outline"
		size="sm"
		class={followLogs ? 'border-primary text-primary' : undefined}
		aria-pressed={followLogs}
		onclick={onEnableFollow}
	>
		Follow logs
	</Button>
	<div class="grid min-w-40 gap-2">
		<Label>Verbosity</Label>
		<SettingsSelect
			value={level}
			options={levelOptions}
			disabled={loading || saving}
			onValueChange={onLevelChange}
		/>
	</div>
</div>
