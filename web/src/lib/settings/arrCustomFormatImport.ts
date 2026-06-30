import type { CustomFormatForm, CustomFormatSpec, CustomFormatSpecType } from './types';

type ArrCustomFormat = {
	name?: unknown;
	specifications?: unknown;
};

type ArrSpecification = {
	name?: unknown;
	implementation?: unknown;
	implementationName?: unknown;
	negate?: unknown;
	required?: unknown;
	fields?: unknown;
	value?: unknown;
};

type ArrField = {
	name?: unknown;
	value?: unknown;
};

export function parseArrCustomFormatImport(raw: string): CustomFormatForm[] {
	const parsed = JSON.parse(raw) as unknown;
	const items = arrItems(parsed);
	if (items.length === 0) {
		throw new Error('No custom formats found in JSON');
	}
	return items.map(arrCustomFormatToForm);
}

function arrItems(value: unknown): ArrCustomFormat[] {
	if (Array.isArray(value)) {
		return value.filter(isObject) as ArrCustomFormat[];
	}
	if (!isObject(value)) {
		return [];
	}
	const wrapper = value as Record<string, unknown>;
	for (const key of ['customFormats', 'formats']) {
		if (Array.isArray(wrapper[key])) {
			return (wrapper[key] as unknown[]).filter(isObject) as ArrCustomFormat[];
		}
	}
	return [wrapper as ArrCustomFormat];
}

function arrCustomFormatToForm(format: ArrCustomFormat): CustomFormatForm {
	const name = stringValue(format.name).trim();
	if (!name) {
		throw new Error('Imported custom format is missing a name');
	}
	const specs = Array.isArray(format.specifications)
		? (format.specifications.filter(isObject) as ArrSpecification[])
		: [];
	const includeSpecs: CustomFormatSpec[] = [];
	const excludeSpecs: CustomFormatSpec[] = [];
	for (const spec of specs) {
		const converted = arrSpecToCustomFormatSpec(spec);
		if (!converted) {
			continue;
		}
		if (spec.negate === true) {
			excludeSpecs.push(converted);
		} else {
			includeSpecs.push(converted);
		}
	}
	if (includeSpecs.length === 0 && excludeSpecs.length === 0) {
		throw new Error(`${name} does not contain importable specifications`);
	}
	return { name, includeSpecs, excludeSpecs };
}

function arrSpecToCustomFormatSpec(spec: ArrSpecification): CustomFormatSpec | undefined {
	const name = stringValue(spec.name).trim() || stringValue(spec.implementationName).trim();
	const value = specValue(spec);
	if (!name || !value) {
		return undefined;
	}
	return {
		id: slug(`${name}-${value}`).slice(0, 80) || `spec-${Date.now()}`,
		name: name.slice(0, 120),
		type: specType(spec),
		value: value.slice(0, 500),
		required: spec.required !== false
	};
}

function specType(spec: ArrSpecification): CustomFormatSpecType {
	const implementation = `${stringValue(spec.implementation)} ${stringValue(spec.implementationName)}`;
	const normalized = implementation.toLowerCase();
	if (normalized.includes('source')) return 'source';
	if (normalized.includes('resolution')) return 'resolution';
	if (normalized.includes('quality')) return 'quality';
	if (normalized.includes('videocodec')) return 'videoCodec';
	if (normalized.includes('audiocodec')) return 'audioCodec';
	if (normalized.includes('releasegroup')) return 'releaseGroup';
	if (normalized.includes('releasetype')) return 'releaseType';
	if (normalized.includes('edition')) return 'edition';
	if (normalized.includes('indexerflag')) return 'indexerFlag';
	if (normalized.includes('language')) return 'language';
	return 'releaseTitle';
}

function specValue(spec: ArrSpecification) {
	const direct = stringValue(spec.value).trim();
	if (direct) {
		return direct;
	}
	const fields = Array.isArray(spec.fields) ? (spec.fields.filter(isObject) as ArrField[]) : [];
	const preferred = [
		'value',
		'pattern',
		'regex',
		'term',
		'source',
		'quality',
		'qualityModifier',
		'resolution',
		'language',
		'codec'
	];
	for (const name of preferred) {
		const field = fields.find((item) => stringValue(item.name) === name);
		const value = fieldValue(field);
		if (value) {
			return value;
		}
	}
	for (const field of fields) {
		const value = fieldValue(field);
		if (value) {
			return value;
		}
	}
	return '';
}

function fieldValue(field: ArrField | undefined) {
	if (!field) {
		return '';
	}
	if (typeof field.value === 'string') {
		return field.value.trim();
	}
	if (typeof field.value === 'number' || typeof field.value === 'boolean') {
		return String(field.value);
	}
	return '';
}

function stringValue(value: unknown) {
	return typeof value === 'string' ? value : '';
}

function isObject(value: unknown): value is Record<string, unknown> {
	return typeof value === 'object' && value !== null;
}

function slug(value: string) {
	return value
		.toLowerCase()
		.replace(/[^a-z0-9]+/g, '-')
		.replace(/^-+|-+$/g, '');
}
