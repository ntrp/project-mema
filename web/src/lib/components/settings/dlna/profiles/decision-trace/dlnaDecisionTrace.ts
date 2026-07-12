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

export type DLNAProfileMatchCandidate = DLNAProfileMatchTraceResponse['candidates'][number];

export interface DLNAProfileMatchView {
	profileId: string;
	profileName: string;
	selectionMethod: string;
	matchReason: string;
	winningRule: string;
	fallbackPath: string;
	score: number;
	headersSummary: string[];
	candidates: DLNAProfileMatchCandidate[];
}

export function basename(path: string) {
	return path.split(/[\\/]/).filter(Boolean).pop() ?? path;
}

export function buildDeliveryTraceSteps(delivery?: DLNADeliveryTraceResponse): DLNATraceStep[] {
	return (delivery?.capabilityTrace ?? []).map((step, index) => ({
		id: `delivery-${index}`,
		stage: 'Delivery decision',
		field: step.field,
		rule: step.rule || 'Rule',
		value: step.value || '—',
		result: step.result
	}));
}

export function buildProfileMatchView(
	profileMatch?: DLNAProfileMatchTraceResponse
): DLNAProfileMatchView | undefined {
	if (!profileMatch) return undefined;
	return {
		profileId: profileMatch.profileId,
		profileName: profileMatch.profileName,
		selectionMethod: matchSourceLabel(profileMatch.matchSource),
		matchReason: profileMatch.matchReason,
		winningRule: profileMatch.winningRule,
		fallbackPath: profileMatch.fallbackPath,
		score: profileMatch.score,
		headersSummary: profileMatch.headersSummary,
		candidates: [...profileMatch.candidates].sort(compareCandidates)
	};
}

function matchSourceLabel(source: string) {
	return (
		{
			manual_uuid: 'Manual device override',
			manual_ip: 'Manual IP override',
			match: 'Automatic profile match',
			sticky_ip: 'Remembered profile for this IP',
			default: 'Generic fallback'
		}[source] ?? source
	);
}

function compareCandidates(a: DLNAProfileMatchCandidate, b: DLNAProfileMatchCandidate) {
	if (a.selected !== b.selected) return a.selected ? -1 : 1;
	if (a.qualified !== b.qualified) return a.qualified ? -1 : 1;
	if (a.priority !== b.priority) return b.priority - a.priority;
	if (a.score !== b.score) return b.score - a.score;
	return a.profileId.localeCompare(b.profileId);
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

export function buildDeliveryTraceSummary(
	delivery?: DLNADeliveryTraceResponse
): DLNATraceSummaryItem[] {
	return [
		{ label: 'Profile name', value: delivery?.profileName ?? '—' },
		{ label: 'Delivery Protocol', value: delivery?.deliveryProtocol ?? '—' },
		{ label: 'Delivery mode', value: deliveryModeDetail(delivery) },
		{ label: 'Video codec', value: codecDetail(delivery, 'videoCodec', delivery?.videoCodec) },
		{ label: 'Audio codec', value: codecDetail(delivery, 'audioCodec', delivery?.audioCodec) }
	];
}

export function filterTraceSteps(steps: DLNATraceStep[], hideFailedSteps = true) {
	return hideFailedSteps ? steps.filter((step) => step.result !== 'fail') : steps;
}
