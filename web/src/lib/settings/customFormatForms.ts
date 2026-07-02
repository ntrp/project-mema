import type { CustomFormat, CustomFormatForm, CustomFormatRequest } from './types';

export function emptyCustomFormatForm(): CustomFormatForm {
	return {
		name: '',
		includeInRenameTemplate: false,
		includeSpecs: [],
		excludeSpecs: []
	};
}

export function customFormatFormFromFormat(format: CustomFormat): CustomFormatForm {
	return {
		id: format.id,
		name: format.name,
		includeInRenameTemplate: format.includeInRenameTemplate,
		includeSpecs: format.includeSpecs.map((spec) => ({ ...spec })),
		excludeSpecs: format.excludeSpecs.map((spec) => ({ ...spec }))
	};
}

export function normalizeCustomFormatForm(form: CustomFormatForm): CustomFormatRequest {
	return {
		name: form.name.trim(),
		includeInRenameTemplate: form.includeInRenameTemplate,
		includeSpecs: normalizeCustomFormatSpecs(form.includeSpecs),
		excludeSpecs: normalizeCustomFormatSpecs(form.excludeSpecs)
	};
}

function normalizeCustomFormatSpecs(specs: CustomFormatRequest['includeSpecs']) {
	return specs
		.map((spec) => ({
			id: spec.id.trim(),
			name: spec.name.trim(),
			type: spec.type,
			value: spec.value.trim(),
			required: spec.required
		}))
		.filter((spec) => spec.id !== '' && spec.name !== '' && spec.value !== '');
}
