import { describe, expect, it } from 'vitest';

import {
	isCandidateVisualState,
	isTargetSatisfactionState,
	TARGET_CANDIDATE_VISUAL_STATES,
	TARGET_SATISFACTION_STATES
} from './targetSatisfactionContract';

describe('target satisfaction contract', () => {
	it('keeps target states distinct from legacy availability states', () => {
		expect(TARGET_SATISFACTION_STATES).not.toContain('available');
		expect(TARGET_SATISFACTION_STATES).not.toContain('unmanaged');
		expect(isTargetSatisfactionState('satisfied')).toBe(true);
		expect(isTargetSatisfactionState('available')).toBe(false);
	});

	it('keeps candidate visual states separate from target states', () => {
		expect(TARGET_CANDIDATE_VISUAL_STATES).toContain('unwanted');
		expect(isCandidateVisualState('missing_placeholder')).toBe(true);
		expect(isCandidateVisualState('failed')).toBe(false);
	});
});
