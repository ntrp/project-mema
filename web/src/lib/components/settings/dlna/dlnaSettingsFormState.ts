import type { DLNASettingsRequest } from '$lib/settings/types';

const defaultAllowedCidrs = ['127.0.0.1/32', '::1/128'];

export function createDLNASettingsForm(settings?: DLNASettingsRequest): DLNASettingsRequest {
	if (!settings) {
		return {
			enabled: false,
			friendlyName: 'Mema',
			interfaces: [],
			allowedCidrs: [...defaultAllowedCidrs],
			announceIntervalSeconds: 1800,
			transcodeEnabled: true,
			thumbnailsEnabled: true,
			subtitlesEnabled: true,
			defaultRendererProfile: 'generic'
		};
	}

	return {
		enabled: settings.enabled,
		friendlyName: settings.friendlyName,
		interfaces: [...settings.interfaces],
		allowedCidrs: [...settings.allowedCidrs],
		announceIntervalSeconds: settings.announceIntervalSeconds,
		transcodeEnabled: settings.transcodeEnabled,
		thumbnailsEnabled: settings.thumbnailsEnabled,
		subtitlesEnabled: settings.subtitlesEnabled,
		defaultRendererProfile: settings.defaultRendererProfile
	};
}

export function allowedCidrsText(allowedCidrs: string[]) {
	return allowedCidrs.join('\n');
}
