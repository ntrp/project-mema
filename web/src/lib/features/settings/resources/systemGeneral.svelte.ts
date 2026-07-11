import { createMutation, createQuery, useQueryClient } from '@tanstack/svelte-query';
import {
	getSystemEventSettings,
	updateSystemEventSettings
} from '$lib/components/settings/system/events/api';
import {
	getSystemLogLevel,
	getSystemLogFileSettings,
	listSystemLogFiles,
	updateSystemLogLevel,
	updateSystemLogFileSettings
} from '$lib/components/settings/system/logs/api';

export const systemSettingsKeys = {
	events: () => ['settings', 'system', 'events'] as const,
	logFile: () => ['settings', 'system', 'log-file'] as const,
	logFiles: () => ['settings', 'system', 'log-files'] as const,
	logLevel: () => ['settings', 'system', 'log-level'] as const
};

function resource<T, V>(
	key: readonly string[],
	queryFn: () => Promise<T>,
	mutationFn: (value: V) => Promise<T>
) {
	const client = useQueryClient();
	return {
		query: createQuery(() => ({ queryKey: key, queryFn })),
		save: createMutation(() => ({
			mutationFn,
			onSuccess: (data) => client.setQueryData(key, data)
		}))
	};
}

export const createEventSettingsResource = () =>
	resource(systemSettingsKeys.events(), getSystemEventSettings, updateSystemEventSettings);

export const createLogFileSettingsResource = () =>
	resource(systemSettingsKeys.logFile(), getSystemLogFileSettings, updateSystemLogFileSettings);

export const createLogFilesQuery = () =>
	createQuery(() => ({ queryKey: systemSettingsKeys.logFiles(), queryFn: listSystemLogFiles }));

export const createLogLevelResource = () =>
	resource(systemSettingsKeys.logLevel(), getSystemLogLevel, updateSystemLogLevel);
