import type { CustomFormatSpec } from '$lib/settings/types';

export type CustomFormatSpecTone = 'mandatory' | 'negated' | 'other';

export function customFormatSpecTone(
	spec: CustomFormatSpec,
	negated = false
): CustomFormatSpecTone {
	if (negated) {
		return 'negated';
	}
	return spec.required ? 'mandatory' : 'other';
}

export function customFormatSpecToneClass(tone: CustomFormatSpecTone) {
	switch (tone) {
		case 'mandatory':
			return 'bg-emerald-600 text-white';
		case 'negated':
			return 'bg-destructive text-destructive-foreground';
		case 'other':
			return 'bg-secondary text-secondary-foreground';
	}
}
