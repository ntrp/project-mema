import type { components } from '$lib/api/generated/schema';

export type TargetSatisfactionState = components['schemas']['TargetSatisfactionState'];
export type TargetSatisfactionType = components['schemas']['TargetSatisfactionType'];
export type TargetCandidateVisualState = components['schemas']['TargetCandidateVisualState'];

export const TARGET_SATISFACTION_STATES = [
	'missing',
	'partial',
	'pending',
	'satisfied',
	'upgradeable',
	'blocked',
	'failed'
] as const satisfies TargetSatisfactionState[];

export const TARGET_CANDIDATE_VISUAL_STATES = [
	'matching',
	'partial',
	'unwanted',
	'pending_operation',
	'missing_placeholder'
] as const satisfies TargetCandidateVisualState[];

export function isTargetSatisfactionState(value: string): value is TargetSatisfactionState {
	return TARGET_SATISFACTION_STATES.includes(value as TargetSatisfactionState);
}

export function isCandidateVisualState(value: string): value is TargetCandidateVisualState {
	return TARGET_CANDIDATE_VISUAL_STATES.includes(value as TargetCandidateVisualState);
}
