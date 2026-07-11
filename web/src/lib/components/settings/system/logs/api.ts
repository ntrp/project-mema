import { client } from '$lib/api/client';
import type {
	SystemLogFile,
	SystemLogFileSettings,
	SystemLogFileSettingsRequest,
	SystemLogLevel,
	SystemLogLevelResponse,
	SystemStatusResponse
} from '$lib/settings/types';

function required<T>(data: T | undefined, message: string): T {
	if (!data) throw new Error(message);
	return data;
}

export async function getSystemLogLevel(): Promise<SystemLogLevelResponse> {
	const { data, error } = await client.GET('/system/log-level');
	if (error) throw new Error(error.message);
	return required(data, 'Log level request did not return a result');
}

export async function updateSystemLogLevel(level: SystemLogLevel) {
	const { data, error } = await client.PUT('/system/log-level', { body: { level } });
	if (error) throw new Error(error.message);
	return required(data, 'Log level update did not return a result');
}

export async function getSystemStatus(): Promise<SystemStatusResponse> {
	const { data, error } = await client.GET('/system/status');
	if (error) throw new Error(error.message);
	return required(data, 'System status request did not return a result');
}

export async function getSystemLogFileSettings(): Promise<SystemLogFileSettings> {
	const { data, error } = await client.GET('/system/log-file-settings');
	if (error) throw new Error(error.message);
	return required(data, 'Log file settings request did not return a result');
}

export async function updateSystemLogFileSettings(request: SystemLogFileSettingsRequest) {
	const { data, error } = await client.PUT('/system/log-file-settings', { body: request });
	if (error) throw new Error(error.message);
	return required(data, 'Log file settings update did not return a result');
}

export async function listSystemLogFiles(): Promise<SystemLogFile[]> {
	const { data, error } = await client.GET('/system/log-files');
	if (error) throw new Error(error.message);
	return data?.files ?? [];
}

export async function downloadSystemLogFile(name: string) {
	const response = await globalThis.fetch(
		`/api/system/log-files/${encodeURIComponent(name)}/download`,
		{
			credentials: 'include'
		}
	);
	if (!response.ok) throw new Error('Could not download log file');
	return response.blob();
}
