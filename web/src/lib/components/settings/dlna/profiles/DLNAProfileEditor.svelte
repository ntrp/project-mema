<script lang="ts">
	import CircleHelpIcon from '@lucide/svelte/icons/circle-help';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import { Switch } from '$lib/components/ui/switch';
	import { Textarea } from '$lib/components/ui/textarea';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import type { DLNARendererProfile } from '$lib/settings/types';
	import type { DLNAProfileForm } from './dlnaProfileForms';
	import { jsonSectionKeys } from './dlnaProfileForms';

	interface Props {
		profile?: DLNARendererProfile;
		form?: DLNAProfileForm;
		saving?: boolean;
		errorMessage?: string;
		onSave: () => void | Promise<void>;
		onCreate: () => void;
		onImport: () => void;
	}

	let {
		profile,
		form = $bindable(),
		saving = false,
		errorMessage = '',
		onSave,
		onCreate,
		onImport
	}: Props = $props();

	const jsonLabels: Record<string, string> = {
		matchRules: 'Matching',
		capabilityRules: 'Direct play',
		deliverySettings: 'Delivery',
		dlnaFlags: 'DLNA flags',
		subtitleRules: 'Subtitles',
		artworkRules: 'Artwork',
		metadataRules: 'Metadata',
		quirks: 'Quirks'
	};
</script>

<section class="grid gap-4" aria-label="DLNA profile editor">
	<div class="flex flex-wrap items-center justify-between gap-3">
		<div class="grid gap-1">
			<h3 class="m-0 text-sm font-semibold">{profile ? profile.name : 'No profile selected'}</h3>
			<p class="m-0 text-xs text-muted-foreground">
				{profile
					? `${profile.source} profile, updated ${new Date(profile.updatedAt).toLocaleString()}`
					: 'Select a profile to edit'}
			</p>
		</div>
		<div class="flex gap-2">
			<Button type="button" variant="outline" size="sm" onclick={onImport}>Import</Button>
			<Button type="button" variant="outline" size="sm" onclick={onCreate}>New custom</Button>
		</div>
	</div>

	{#if errorMessage}
		<p class="m-0 text-sm font-medium text-destructive">{errorMessage}</p>
	{/if}

	{#if form}
		<form
			class="grid gap-4"
			onsubmit={(event) => {
				event.preventDefault();
				void onSave();
			}}
		>
			<div class="grid gap-4 sm:grid-cols-4">
				<div class="space-y-2">
					<Label for="dlna-profile-id">ID</Label>
					<Input id="dlna-profile-id" bind:value={form.id} disabled={Boolean(profile)} required />
				</div>
				<div class="space-y-2 sm:col-span-2">
					<Label for="dlna-profile-name">Name</Label>
					<Input id="dlna-profile-name" bind:value={form.name} required />
				</div>
				<label class="flex items-center gap-3 pt-7">
					<Switch bind:checked={form.enabled} />
					<span class="text-sm font-medium">Enabled</span>
				</label>
				<div class="space-y-2">
					<Label for="dlna-profile-vendor">Family</Label>
					<Input id="dlna-profile-vendor" bind:value={form.vendor} />
				</div>
				<div class="space-y-2">
					<Label for="dlna-profile-class">Device class</Label>
					<Input id="dlna-profile-class" bind:value={form.deviceClass} required />
				</div>
				<div class="space-y-2">
					<Label for="dlna-profile-priority">Priority</Label>
					<Input id="dlna-profile-priority" type="number" bind:value={form.priority} />
				</div>
				<div class="space-y-2">
					<Label for="dlna-profile-icon">Icon key</Label>
					<Input id="dlna-profile-icon" bind:value={form.iconKey} />
				</div>
			</div>
			<div class="space-y-2">
				<Label for="dlna-profile-notes">Notes</Label>
				<Textarea id="dlna-profile-notes" class="min-h-20" bind:value={form.notes} />
			</div>
			<div class="grid gap-4 lg:grid-cols-2">
				{#each jsonSectionKeys as key (key)}
					<div class="space-y-2">
						<div class="flex items-center gap-2">
							<Label for={`dlna-profile-${key}`}>{jsonLabels[key]}</Label>
							<Tooltip.Root>
								<Tooltip.Trigger>
									{#snippet child({ props })}
										<Button
											{...props}
											type="button"
											variant="ghost"
											size="icon-sm"
											aria-label={`${jsonLabels[key]} JSON help`}
										>
											<CircleHelpIcon aria-hidden="true" />
										</Button>
									{/snippet}
								</Tooltip.Trigger>
								<Tooltip.Content>JSON object used by renderer profile evaluator</Tooltip.Content>
							</Tooltip.Root>
						</div>
						<Textarea
							id={`dlna-profile-${key}`}
							class="min-h-36 font-mono text-xs"
							bind:value={form.jsonText[key]}
							spellcheck={false}
						/>
					</div>
				{/each}
			</div>
			<div class="flex justify-end">
				<Button type="submit" disabled={saving}
					>{saving ? 'Saving' : profile ? 'Save profile' : 'Create profile'}</Button
				>
			</div>
		</form>
	{/if}
</section>
