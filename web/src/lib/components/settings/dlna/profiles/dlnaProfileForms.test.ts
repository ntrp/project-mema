import { describe, expect, it } from 'vitest';
import { renderWithTooltip } from '$lib/components/rendered/renderHelpers';
import type { DLNARendererProfile } from '$lib/settings/types';
import DLNADeviceOverrideTable from './DLNADeviceOverrideTable.svelte';
import DLNAProfileTable from './DLNAProfileTable.svelte';
import { formToRequest, importProfileText, profileToForm } from './dlnaProfileForms';

describe('DLNA profile settings helpers (SCN-SETTINGS-025)', () => {
	it('round-trips seeded profile JSON fields for editing', () => {
		const form = profileToForm(sampleProfile);
		form.jsonText.capabilityRules = '{"containers":["mkv"]}';

		expect(formToRequest(form)).toMatchObject({
			name: 'LG webOS',
			capabilityRules: { containers: ['mkv'] }
		});
	});

	it('rejects imported profiles without required identity', () => {
		expect(() => importProfileText('{"name":"Broken"}')).toThrow('missing id');
	});
});

describe('DLNA device profile settings UI (SCN-SETTINGS-025)', () => {
	it('renders profile and override tables', () => {
		const profileTable = renderWithTooltip(DLNAProfileTable, {
			profiles: [sampleProfile],
			search: '',
			selectedId: sampleProfile.id,
			onSearch: () => {},
			onSelect: () => {},
			onClone: () => {},
			onReset: () => {},
			onExport: () => {}
		});
		const devices = renderWithTooltip(DLNADeviceOverrideTable, {
			devices: [
				{
					ip: '192.168.1.55',
					userAgent: 'LG TV',
					profileId: sampleProfile.id,
					lastSoapAction: 'Browse',
					lastSeen: '2026-07-08T08:00:00Z'
				}
			],
			overrides: [],
			profiles: [sampleProfile],
			overrideForm: {
				ipAddress: '',
				rendererUuid: '',
				profileId: sampleProfile.id,
				displayName: '',
				allowed: true,
				deliveryPolicyOverrides: {},
				notes: ''
			},
			overrideJsonText: '{}',
			onOverrideJson: () => {},
			onSave: () => {},
			onDelete: () => {},
			onQuickAssign: () => {}
		});

		expect(profileTable.body).toContain('LG webOS');
		expect(profileTable.body).toContain('Customized');
		expect(devices.body).toContain('192.168.1.55');
		expect(devices.body).toContain('Manual override');
	});
});

const sampleProfile: DLNARendererProfile = {
	id: 'lg-webos',
	name: 'LG webOS',
	vendor: 'LG',
	deviceClass: 'MediaRenderer',
	enabled: true,
	priority: 120,
	iconKey: 'tv',
	notes: 'seeded',
	matchRules: { headers: ['LG'] },
	capabilityRules: { containers: ['mp4'] },
	deliverySettings: {},
	dlnaFlags: {},
	subtitleRules: {},
	artworkRules: {},
	metadataRules: {},
	quirks: {},
	source: 'mema_seed',
	sourceVersion: 1,
	customized: false,
	createdAt: '2026-07-08T08:00:00Z',
	updatedAt: '2026-07-08T08:00:00Z'
};
