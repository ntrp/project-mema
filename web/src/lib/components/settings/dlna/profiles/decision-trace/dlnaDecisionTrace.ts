import type { DLNADeliveryTraceResponse, DLNAProfileMatchTraceResponse } from '$lib/settings/types';

export interface DLNATraceStep {
	id: string;
	stage: string;
	field: string;
	rule: string;
	value: string;
	score?: number;
	result: string;
}

export interface DLNATraceSummaryItem {
	label: string;
	value: string;
}

export function basename(path: string) {
	return path.split(/[\\/]/).filter(Boolean).pop() ?? path;
}

export function buildTraceSteps(
	profileMatch?: DLNAProfileMatchTraceResponse,
	delivery?: DLNADeliveryTraceResponse
): DLNATraceStep[] {
	const profileSteps = profileMatch
		? [
				{
					id: 'profile-selected',
					stage: 'Profile matching',
					field: 'selected profile',
					rule: profileMatch.matchReason || 'Profile selection',
					value: `${profileMatch.profileName} (${profileMatch.profileId})`,
					score: profileMatch.score,
					result: 'pass'
				},
				...profileMatch.ruleTrace.map((step, index) => ({
					id: `profile-${index}`,
					stage: 'Profile matching',
					field: step.field,
					rule: `${step.profileName} (${step.profileId}): ${step.rule || 'Rule'}`,
					value: step.value || '—',
					score: step.score,
					result: step.result
				}))
			]
		: [];
	const deliverySteps = (delivery?.capabilityTrace ?? []).map((step, index) => ({
		id: `delivery-${index}`,
		stage: 'Delivery decision',
		field: step.field,
		rule: step.rule || 'Rule',
		value: step.value || '—',
		result: step.result
	}));

	return [...profileSteps, ...deliverySteps];
}

function tracedValue(delivery: DLNADeliveryTraceResponse | undefined, field: string) {
	return delivery?.capabilityTrace.find((step) => step.field === field)?.value;
}

function deliveryModeDetail(delivery?: DLNADeliveryTraceResponse) {
	if (!delivery) return '—';
	if (delivery.mode !== 'remux') return delivery.mode;
	return `remux (${tracedValue(delivery, 'container') ?? '—'} -> mpegts)`;
}

function codecDetail(
	delivery: DLNADeliveryTraceResponse | undefined,
	field: 'videoCodec' | 'audioCodec',
	outputCodec: string | undefined
) {
	if (!delivery) return '—';
	const sourceCodec = tracedValue(delivery, field)?.trim();
	const effectiveOutputCodec = outputCodec || (delivery.mode === 'direct' ? 'copy' : undefined);
	if (!effectiveOutputCodec) return sourceCodec || '—';
	if (!sourceCodec)
		return effectiveOutputCodec === 'copy' ? 'copy' : `transcoding (${effectiveOutputCodec})`;
	return effectiveOutputCodec === 'copy'
		? `copy (${sourceCodec})`
		: `transcoding (${sourceCodec} -> ${effectiveOutputCodec})`;
}

export function buildTraceSummary(
	profileMatch?: DLNAProfileMatchTraceResponse,
	delivery?: DLNADeliveryTraceResponse
): DLNATraceSummaryItem[] {
	return [
		{ label: 'Profile name', value: delivery?.profileName ?? profileMatch?.profileName ?? '—' },
		{ label: 'Delivery Protocol', value: delivery?.deliveryProtocol ?? '—' },
		{ label: 'Delivery mode', value: deliveryModeDetail(delivery) },
		{ label: 'Video codec', value: codecDetail(delivery, 'videoCodec', delivery?.videoCodec) },
		{ label: 'Audio codec', value: codecDetail(delivery, 'audioCodec', delivery?.audioCodec) },
		{ label: 'Match score', value: profileMatch ? String(profileMatch.score) : '—' },
		{ label: 'Match reason', value: profileMatch?.matchReason ?? '—' },
		{ label: 'Winning rule', value: profileMatch?.winningRule ?? '—' }
	];
}

export function filterTraceSteps(steps: DLNATraceStep[], hideFailedSteps = true) {
	return hideFailedSteps ? steps.filter((step) => step.result !== 'fail') : steps;
}
