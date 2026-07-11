import { describe, expect, it, vi } from 'vitest';

const generated = vi.hoisted(() => ({
	listMediaItems: vi.fn(),
	listMediaRequests: vi.fn()
}));

vi.mock('$lib/api/generated/tanstack', () => generated);

import { listMediaItems, listMediaRequests } from './api';

describe('library API boundary', () => {
	it('exposes the generated operations owned by the library feature', () => {
		expect(listMediaItems).toBe(generated.listMediaItems);
		expect(listMediaRequests).toBe(generated.listMediaRequests);
	});
});
