import { beforeEach, describe, expect, it, vi } from 'vitest';

const mocks = vi.hoisted(() => ({
	client: { POST: vi.fn(), PUT: vi.fn() },
	bodies: { metadata: {}, subtitle: {}, profile: {} }
}));
vi.mock('$lib/api/client', () => ({ client: mocks.client }));
vi.mock('../forms', () => ({
	normalizeMetadataProviderForm: () => mocks.bodies.metadata,
	normalizeSubtitleProviderForm: () => mocks.bodies.subtitle,
	normalizeMediaProfileForm: () => mocks.bodies.profile
}));

import {
	saveMediaProfile,
	saveMetadataProvider,
	saveSubtitleProvider,
	testMetadataProvider,
	testSubtitleProvider
} from './providersProfiles';

describe('provider and profile settings domain', () => {
	beforeEach(() => vi.clearAllMocks());

	it('creates and updates each normalized entity', async () => {
		mocks.client.POST.mockResolvedValue({});
		mocks.client.PUT.mockResolvedValue({});
		for (const save of [saveMetadataProvider, saveSubtitleProvider, saveMediaProfile]) {
			await save({} as never);
			await save({ id: 'entity-1' } as never);
		}
		expect(mocks.client.POST).toHaveBeenCalledTimes(3);
		expect(mocks.client.PUT).toHaveBeenCalledTimes(3);
	});

	it('returns provider test results', async () => {
		mocks.client.POST.mockResolvedValue({ data: { success: true } });
		await expect(testMetadataProvider('metadata')).resolves.toEqual({ success: true });
		await expect(testSubtitleProvider('subtitle')).resolves.toEqual({ success: true });
	});

	it('surfaces create, update, and test errors', async () => {
		mocks.client.POST.mockResolvedValue({ error: { message: 'post failed' } });
		for (const command of [
			() => saveMetadataProvider({} as never),
			() => saveSubtitleProvider({} as never),
			() => saveMediaProfile({} as never),
			() => testMetadataProvider('id'),
			() => testSubtitleProvider('id')
		])
			await expect(command()).rejects.toThrow('post failed');

		mocks.client.PUT.mockResolvedValue({ error: { message: 'put failed' } });
		for (const save of [saveMetadataProvider, saveSubtitleProvider, saveMediaProfile])
			await expect(save({ id: 'id' } as never)).rejects.toThrow('put failed');
	});

	it('rejects missing provider test results', async () => {
		mocks.client.POST.mockResolvedValue({});
		await expect(testMetadataProvider('id')).rejects.toThrow('did not return');
		await expect(testSubtitleProvider('id')).rejects.toThrow('did not return');
	});
});
