import { describe, expect, it } from 'vitest';

import type { DLNASettingsRequest } from '$lib/settings/types';
import { allowedCidrsText, createDLNASettingsForm } from './dlnaSettingsFormState';

describe('DLNA settings form state helpers', () => {
	it('creates fresh default settings copies', () => {
		const form = createDLNASettingsForm();

		expect(form).toEqual({
			enabled: false,
			friendlyName: 'Mema',
			interfaces: [],
			allowedCidrs: ['127.0.0.1/32', '::1/128'],
			announceIntervalSeconds: 1800,
			transcodeEnabled: true,
			thumbnailsEnabled: true,
			subtitlesEnabled: true,
			defaultRendererProfile: 'generic'
		});

		form.interfaces.push('eth0');
		form.allowedCidrs.push('10.0.0.0/8');

		expect(createDLNASettingsForm()).toEqual({
			enabled: false,
			friendlyName: 'Mema',
			interfaces: [],
			allowedCidrs: ['127.0.0.1/32', '::1/128'],
			announceIntervalSeconds: 1800,
			transcodeEnabled: true,
			thumbnailsEnabled: true,
			subtitlesEnabled: true,
			defaultRendererProfile: 'generic'
		});
	});

	it('clones loaded settings and formats allowed CIDRs', () => {
		const settings: DLNASettingsRequest = {
			enabled: true,
			friendlyName: 'Living Room TV',
			interfaces: ['eth0'],
			allowedCidrs: ['192.168.1.0/24'],
			announceIntervalSeconds: 900,
			transcodeEnabled: false,
			thumbnailsEnabled: false,
			subtitlesEnabled: true,
			defaultRendererProfile: 'lg-webos'
		};

		const form = createDLNASettingsForm(settings);

		expect(form).not.toBe(settings);
		expect(form.interfaces).not.toBe(settings.interfaces);
		expect(form.allowedCidrs).not.toBe(settings.allowedCidrs);
		expect(form).toEqual(settings);
		expect(allowedCidrsText(form.allowedCidrs)).toBe('192.168.1.0/24');
	});
});
