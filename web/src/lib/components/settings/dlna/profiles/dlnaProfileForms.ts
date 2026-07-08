import type {
	DLNARendererDeviceOverrideRequest,
	DLNARendererProfile,
	DLNARendererProfileCreateRequest,
	DLNARendererProfileRequest
} from '$lib/settings/types';

export const jsonSectionKeys = [
	'matchRules',
	'capabilityRules',
	'deliverySettings',
	'dlnaFlags',
	'subtitleRules',
	'artworkRules',
	'metadataRules',
	'quirks'
] as const;

export type DLNAJsonSectionKey = (typeof jsonSectionKeys)[number];
export type DLNAProfileJsonText = Record<DLNAJsonSectionKey, string>;

export interface DLNAProfileForm extends DLNARendererProfileRequest {
	id: string;
	jsonText: DLNAProfileJsonText;
}

export const emptyJsonObject = {};

export function profileToForm(profile: DLNARendererProfile): DLNAProfileForm {
	return {
		id: profile.id,
		name: profile.name,
		vendor: profile.vendor,
		deviceClass: profile.deviceClass,
		enabled: profile.enabled,
		priority: profile.priority,
		iconKey: profile.iconKey,
		notes: profile.notes,
		matchRules: profile.matchRules,
		capabilityRules: profile.capabilityRules,
		deliverySettings: profile.deliverySettings,
		dlnaFlags: profile.dlnaFlags,
		subtitleRules: profile.subtitleRules,
		artworkRules: profile.artworkRules,
		metadataRules: profile.metadataRules,
		quirks: profile.quirks,
		jsonText: jsonSectionText(profile)
	};
}

export function formToRequest(form: DLNAProfileForm): DLNARendererProfileRequest {
	const parsed = parseJsonSections(form.jsonText);
	return {
		name: form.name.trim(),
		vendor: form.vendor.trim(),
		deviceClass: form.deviceClass.trim(),
		enabled: form.enabled,
		priority: Number(form.priority),
		iconKey: form.iconKey.trim(),
		notes: form.notes.trim(),
		...parsed
	};
}

export function formToCreateRequest(form: DLNAProfileForm): DLNARendererProfileCreateRequest {
	return {
		id: form.id.trim(),
		...formToRequest(form)
	};
}

export function blankProfileForm(source?: DLNARendererProfile): DLNAProfileForm {
	const base = source ? profileToForm(source) : profileToForm(defaultProfile());
	return {
		...base,
		id: source ? `${source.id}-copy` : '',
		name: source ? `${source.name} copy` : ''
	};
}

export function defaultOverrideRequest(profileId = ''): DLNARendererDeviceOverrideRequest {
	return {
		rendererUuid: '',
		ipAddress: '',
		profileId,
		displayName: '',
		allowed: true,
		deliveryPolicyOverrides: emptyJsonObject,
		notes: ''
	};
}

export function importProfileText(text: string): DLNARendererProfileCreateRequest {
	const value = JSON.parse(text) as DLNARendererProfileCreateRequest;
	for (const key of ['id', 'name', 'vendor', 'deviceClass']) {
		if (typeof value[key as keyof DLNARendererProfileCreateRequest] !== 'string') {
			throw new Error(`Imported profile missing ${key}`);
		}
	}
	return value;
}

export function profileExportText(profile: DLNARendererProfile) {
	return JSON.stringify(profile, null, 2);
}

export function downloadProfileJson(profile: DLNARendererProfile) {
	const blob = new Blob([profileExportText(profile)], { type: 'application/json' });
	const url = URL.createObjectURL(blob);
	const link = document.createElement('a');
	link.href = url;
	link.download = `${profile.id}.json`;
	link.click();
	URL.revokeObjectURL(url);
}

function jsonSectionText(profile: DLNARendererProfileRequest): DLNAProfileJsonText {
	return Object.fromEntries(
		jsonSectionKeys.map((key) => [key, JSON.stringify(profile[key] ?? emptyJsonObject, null, 2)])
	) as DLNAProfileJsonText;
}

function parseJsonSections(
	text: DLNAProfileJsonText
): Pick<DLNARendererProfileRequest, DLNAJsonSectionKey> {
	return Object.fromEntries(
		jsonSectionKeys.map((key) => [key, parseJsonObject(text[key], key)])
	) as Pick<DLNARendererProfileRequest, DLNAJsonSectionKey>;
}

function parseJsonObject(text: string, label: string) {
	const parsed = JSON.parse(text || '{}');
	if (!parsed || Array.isArray(parsed) || typeof parsed !== 'object') {
		throw new Error(`${label} must be a JSON object`);
	}
	return parsed as Record<string, unknown>;
}

function defaultProfile(): DLNARendererProfile {
	return {
		id: '',
		name: '',
		vendor: '',
		deviceClass: 'MediaRenderer',
		enabled: true,
		priority: 100,
		iconKey: '',
		notes: '',
		matchRules: emptyJsonObject,
		capabilityRules: emptyJsonObject,
		deliverySettings: emptyJsonObject,
		dlnaFlags: emptyJsonObject,
		subtitleRules: emptyJsonObject,
		artworkRules: emptyJsonObject,
		metadataRules: emptyJsonObject,
		quirks: emptyJsonObject,
		source: 'user',
		sourceVersion: 1,
		customized: false,
		createdAt: '',
		updatedAt: ''
	};
}
