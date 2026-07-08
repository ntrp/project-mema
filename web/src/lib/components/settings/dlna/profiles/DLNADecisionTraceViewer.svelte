<script lang="ts">
	import RefreshCwIcon from '@lucide/svelte/icons/refresh-cw';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import { Textarea } from '$lib/components/ui/textarea';
	import {
		traceDLNADeliveryDecision,
		traceDLNAProfileMatch
	} from '$lib/settings/dlnaProfilesApi';
	import type {
		DLNAClientDiagnostic,
		DLNARendererDeviceOverride,
		DLNARendererProfile
	} from '$lib/settings/types';

	interface Props {
		devices: DLNAClientDiagnostic[];
		overrides: DLNARendererDeviceOverride[];
		profiles: DLNARendererProfile[];
		selectedIp: string;
		mediaPath: string;
		onSelectedIp: (_value: string) => void;
		onMediaPath: (_value: string) => void;
	}

	let { devices, overrides, profiles, selectedIp, mediaPath, onSelectedIp, onMediaPath }: Props =
		$props();

	let loading = $state(false);
	let errorMessage = $state('');
	let traceText = $state('Select a device, enter a media path, then run trace.');

	const selectedDevice = $derived(devices.find((device) => device.ip === selectedIp));
	const selectedOverride = $derived(
		overrides.find((override) => override.ipAddress === selectedIp)
	);
	const selectedProfile = $derived(
		profiles.find(
			(profile) => profile.id === (selectedOverride?.profileId ?? selectedDevice?.profileId)
		)
	);

	async function refreshTrace() {
		loading = true;
		errorMessage = '';
		try {
			const profileMatch = await traceDLNAProfileMatch({
				deviceIp: selectedIp || undefined,
				rendererUuid: selectedDevice?.rendererUuid || undefined,
				friendlyName: selectedDevice?.friendlyName || undefined,
				userAgent: selectedDevice?.userAgent || undefined
			});
			const deliveryDecision = await traceDLNADeliveryDecision({
				deviceIp: selectedIp || undefined,
				profileId: selectedProfile?.id,
				mediaPath: mediaPath || undefined,
				objectId: selectedDevice?.lastObjectId || undefined,
				resourceId: selectedDevice?.lastResourceId || undefined,
				streamMode: selectedDevice?.lastStreamMode || undefined,
				userAgent: selectedDevice?.userAgent || undefined
			});
			traceText = JSON.stringify({ profileMatch, deliveryDecision }, null, 2);
		} catch (error) {
			errorMessage = error instanceof Error ? error.message : 'Could not run DLNA trace';
		} finally {
			loading = false;
		}
	}
</script>

<section class="grid gap-3" aria-label="DLNA decision trace">
	<div class="flex items-center justify-between gap-3">
		<h3 class="m-0 text-sm font-semibold">Decision trace</h3>
		<Button type="button" variant="secondary" size="sm" disabled={loading} onclick={refreshTrace}>
			<RefreshCwIcon class={loading ? 'animate-spin' : ''} />
			Run trace
		</Button>
	</div>
	{#if errorMessage}
		<p class="m-0 text-sm font-medium text-destructive">{errorMessage}</p>
	{/if}
	<div class="grid gap-4 sm:grid-cols-2">
		<div class="space-y-2">
			<Label for="dlna-trace-device">Device</Label>
			<select
				id="dlna-trace-device"
				class="h-9 w-full rounded-md border border-input bg-background px-2 text-sm"
				value={selectedIp}
				onchange={(event) => onSelectedIp(event.currentTarget.value)}
			>
				<option value="">Choose device</option>
				{#each devices as device (device.ip)}
					<option value={device.ip}>{device.ip} - {device.userAgent || device.profileId}</option>
				{/each}
			</select>
		</div>
		<div class="space-y-2">
			<Label for="dlna-trace-media">Media file</Label>
			<Input
				id="dlna-trace-media"
				value={mediaPath}
				oninput={(event) => onMediaPath(event.currentTarget.value)}
				placeholder="/path/to/movie.mkv"
			/>
		</div>
	</div>
	<Textarea class="min-h-52 font-mono text-xs" value={traceText} readonly spellcheck={false} />
</section>
