<script lang="ts">
	import RefreshCwIcon from '@lucide/svelte/icons/refresh-cw';
	import SettingsFormModal from '$lib/components/settings/shared/SettingsFormModal.svelte';
	import { Checkbox } from '$lib/components/ui/checkbox';
	import { Button } from '$lib/components/ui/button';
	import { Label } from '$lib/components/ui/label';
	import { createQuery, useQueryClient } from '@tanstack/svelte-query';
	import { getMediaItemFilePreviewInfo } from '$lib/features/library/filesApi';
	import { traceDLNADeliveryDecision, traceDLNAProfileMatch } from '$lib/settings/dlnaProfilesApi';
	import type {
		DLNAClientDiagnostic,
		DLNARendererProfile,
		MediaFilePreviewInfo
	} from '$lib/settings/types';
	import {
		buildDeliveryTraceSteps,
		buildDeliveryTraceSummary,
		buildProfileMatchView
	} from './dlnaDecisionTrace';
	import { probeFromMediaPreview } from './dlnaMediaProbe';
	import DLNAMediaPicker from './DLNAMediaPicker.svelte';
	import DLNAProfileMatchResults from './DLNAProfileMatchResults.svelte';
	import DLNATraceResults from './DLNATraceResults.svelte';

	interface Props {
		open: boolean;
		devices: DLNAClientDiagnostic[];
		profiles: DLNARendererProfile[];
		selectedIp: string;
		mediaPath: string;
		onClose: () => void;
		onSelectedIp: (_value: string) => void;
		onMediaPath: (_value: string) => void;
	}

	let {
		open,
		devices,
		profiles,
		selectedIp,
		mediaPath,
		onClose,
		onSelectedIp,
		onMediaPath
	}: Props = $props();

	let hideFailedSteps = $state(true);
	let selectedMediaItemId = $state('');
	const queryClient = useQueryClient();

	const selectedDevice = $derived(devices.find((device) => device.ip === selectedIp));
	const selectedProfile = $derived(
		profiles.find((profile) => profile.id === selectedDevice?.profileId)
	);
	const traceHeaders = $derived(headersFromSummary(selectedDevice?.headersSummary ?? []));
	const profileMatch = createQuery(() => ({
		queryKey: [
			'settings',
			'dlna',
			'profile-trace',
			selectedIp,
			selectedDevice?.rendererUuid,
			selectedDevice?.headersSummary
		],
		queryFn: () =>
			traceDLNAProfileMatch({
				deviceIp: selectedIp || undefined,
				rendererUuid: selectedDevice?.rendererUuid || undefined,
				friendlyName: selectedDevice?.friendlyName || undefined,
				userAgent: selectedDevice?.userAgent || undefined,
				headers: traceHeaders
			}),
		enabled: false
	}));
	const delivery = createQuery(() => ({
		queryKey: [
			'settings',
			'dlna',
			'delivery-trace',
			selectedIp,
			selectedProfile?.id,
			selectedMediaItemId,
			mediaPath
		],
		queryFn: async () => {
			const previewInfo = await loadSelectedMediaPreview().catch(() => undefined);
			const probe = probeFromMediaPreview(previewInfo);
			return traceDLNADeliveryDecision({
				deviceIp: selectedIp || undefined,
				profileId: selectedProfile?.id,
				mediaPath: mediaPath || undefined,
				objectId: selectedDevice?.lastObjectId || undefined,
				resourceId: selectedDevice?.lastResourceId || undefined,
				streamMode: selectedDevice?.lastStreamMode || undefined,
				userAgent: selectedDevice?.userAgent || undefined,
				probe
			});
		},
		enabled: false
	}));
	const loading = $derived(profileMatch.isFetching || delivery.isFetching);
	const errorMessage = $derived(profileMatch.error?.message ?? delivery.error?.message ?? '');
	const profileMatchView = $derived(buildProfileMatchView(profileMatch.data));
	const deliverySteps = $derived(buildDeliveryTraceSteps(delivery.data));
	const deliverySummary = $derived(buildDeliveryTraceSummary(delivery.data));

	function headersFromSummary(headers: string[]) {
		return Object.fromEntries(
			headers.flatMap((header) => {
				const separator = header.indexOf(':');
				const name = header.slice(0, separator).trim();
				const value = header.slice(separator + 1).trim();
				return separator > 0 && name && value ? [[name, value]] : [];
			})
		);
	}

	async function loadSelectedMediaPreview(): Promise<MediaFilePreviewInfo | undefined> {
		if (!selectedMediaItemId || !mediaPath) return undefined;
		return getMediaItemFilePreviewInfo(selectedMediaItemId, mediaPath);
	}

	async function refreshTrace() {
		if (!selectedDevice || !mediaPath) return;
		await Promise.all([profileMatch.refetch(), delivery.refetch()]);
	}

	function closeViewer() {
		onClose();
		void Promise.all([
			queryClient.invalidateQueries({ queryKey: ['settings', 'dlna', 'profile-trace'] }),
			queryClient.invalidateQueries({ queryKey: ['settings', 'dlna', 'delivery-trace'] })
		]);
	}
</script>

{#if open}
	<SettingsFormModal
		title="Decision trace"
		onClose={closeViewer}
		modalClass="w-[min(960px,calc(100vw-32px))]"
	>
		<section class="grid gap-4" aria-label="DLNA decision trace">
			{#if errorMessage}
				<p class="m-0 text-sm font-medium text-destructive">{errorMessage}</p>
			{/if}
			<div class="grid gap-4 sm:grid-cols-[minmax(0,1fr)_minmax(0,1fr)_minmax(0,1fr)] sm:items-end">
				<div class="space-y-2 col-span-2">
					<select
						id="dlna-trace-device"
						class="h-9 w-full rounded-md border border-input bg-background px-2 text-sm"
						value={selectedIp}
						onchange={(event) => onSelectedIp(event.currentTarget.value)}
					>
						<option value="">Choose device</option>
						{#each devices as device (device.ip)}
							<option value={device.ip}>{device.ip} - {device.userAgent || device.profileId}</option
							>
						{/each}
					</select>
					<DLNAMediaPicker
						{mediaPath}
						{onMediaPath}
						onSelectedMediaId={(value) => (selectedMediaItemId = value)}
					/>
				</div>
				<div class="flex flex-col gap-3 sm:items-end sm:justify-self-end">
					<div class="flex h-9 items-center justify-end gap-2">
						<Checkbox id="dlna-trace-hide-failed" bind:checked={hideFailedSteps} />
						<Label for="dlna-trace-hide-failed" class="text-sm">Hide failed delivery tests</Label>
					</div>
					<Button
						type="button"
						variant="secondary"
						size="sm"
						disabled={loading || !selectedDevice || !mediaPath}
						onclick={refreshTrace}
						class="w-full sm:w-auto"
					>
						<RefreshCwIcon class={loading ? 'animate-spin' : ''} />
						Run trace
					</Button>
				</div>
			</div>
			<DLNAProfileMatchResults match={profileMatchView} />
			<section class="grid gap-3" aria-label="DLNA delivery decision">
				<h3 class="m-0 text-sm font-semibold">Delivery decision</h3>
				<DLNATraceResults steps={deliverySteps} summary={deliverySummary} {hideFailedSteps} />
			</section>
		</section>
	</SettingsFormModal>
{/if}
