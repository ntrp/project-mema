import { createMutation } from '@tanstack/svelte-query';

export type RunCommandMutation = <T>(command: () => Promise<T>) => Promise<T>;

export function createCommandMutation(): RunCommandMutation {
	const mutation = createMutation(() => ({
		mutationFn: (command: () => Promise<unknown>) => command()
	}));
	return <T>(command: () => Promise<T>) => mutation.mutateAsync(command) as Promise<T>;
}
