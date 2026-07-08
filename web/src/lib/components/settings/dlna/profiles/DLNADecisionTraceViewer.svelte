<script lang="ts">
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import { Textarea } from '$lib/components/ui/textarea';
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

	const selectedDevice = $derived(devices.find((device) => device.ip === selectedIp));
	const selectedOverride = $derived(
		overrides.find((override) => override.ipAddress === selectedIp)
	);
	const selectedProfile = $derived(
		profiles.find(
			(profile) => profile.id === (selectedOverride?.profileId ?? selectedDevice?.profileId)
		)
	);
	const traceText = $derived(
		JSON.stringify(
			{
				deviceIp: selectedIp || null,
				mediaFile: mediaPath || null,
				overrideProfileId: selectedOverride?.profileId ?? null,
				matchedProfileId: selectedDevice?.profileId ?? null,
				effectiveProfileId: selectedProfile?.id ?? null,
				effectiveProfileName: selectedProfile?.name ?? null,
				directPlayRules: selectedProfile?.capabilityRules ?? {},
				deliverySettings: selectedProfile?.deliverySettings ?? {},
				quirks: selectedProfile?.quirks ?? {}
			},
			null,
			2
		)
	);
</script>

<section class="grid gap-3" aria-label="DLNA decision trace">
	<h3 class="m-0 text-sm font-semibold">Decision trace</h3>
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
