import { describe, expect, it, vi } from 'vitest';

const { createMutation } = vi.hoisted(() => ({
	createMutation: vi.fn(
		(options: () => { mutationFn: (command: () => Promise<unknown>) => Promise<unknown> }) => ({
			mutateAsync: (command: () => Promise<unknown>) => options().mutationFn(command)
		})
	)
}));

vi.mock('@tanstack/svelte-query', () => ({ createMutation }));

import { createCommandMutation } from './commandMutation.svelte';

describe('command mutation', () => {
	it('executes server commands through a TanStack mutation', async () => {
		const command = vi.fn().mockResolvedValue({ id: 'saved' });
		await expect(createCommandMutation()(command)).resolves.toEqual({ id: 'saved' });
		expect(createMutation).toHaveBeenCalledOnce();
		expect(command).toHaveBeenCalledOnce();
	});

	it('preserves command failures', async () => {
		const error = new Error('failed');
		await expect(createCommandMutation()(() => Promise.reject(error))).rejects.toBe(error);
	});
});
