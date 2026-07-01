import type { QualitySizeSetting, QualitySizeSettingRequest } from '$lib/settings/types';

export type SliderField = 'minimum' | 'preferred' | 'maximum';

export const sliderMaxGibPerHour = 120;
export const sliderStepGibPerHour = 0.1;
export const sliderHandleGap = 0.1;

export function rowError(quality: QualitySizeSetting) {
	const minimum = quality.minimumSizeMbPerMinute;
	const preferred = quality.preferredSizeMbPerMinute;
	const maximum = quality.maximumSizeMbPerMinute;
	if (minimum < 0) {
		return 'Minimum must be zero or greater';
	}
	if (preferred != null && preferred < minimum) {
		return 'Preferred must be at least minimum';
	}
	if (maximum != null && maximum < minimum) {
		return 'Maximum must be at least minimum';
	}
	if (preferred != null && maximum != null && preferred > maximum) {
		return 'Preferred must be at most maximum';
	}
	return '';
}

export function qualityRequest(quality: QualitySizeSetting): QualitySizeSettingRequest {
	return {
		qualityId: quality.qualityId,
		minimumSizeMbPerMinute: quality.minimumSizeMbPerMinute,
		preferredSizeMbPerMinute: quality.preferredSizeMbPerMinute ?? null,
		maximumSizeMbPerMinute: quality.maximumSizeMbPerMinute ?? null
	};
}

export function mbPerMinuteToGibPerHour(value: number | null | undefined) {
	if (value == null) {
		return sliderMaxGibPerHour;
	}
	return clamp(Math.round(((value * 60) / 1024) * 100) / 100, 0, sliderMaxGibPerHour);
}

export function gibPerHourToMbPerMinute(value: number) {
	return Math.round(((value * 1024) / 60) * 100) / 100;
}

export function sliderValues(quality: QualitySizeSetting) {
	const minimum = mbPerMinuteToGibPerHour(quality.minimumSizeMbPerMinute);
	const maximum = mbPerMinuteToGibPerHour(quality.maximumSizeMbPerMinute);
	const preferred = clamp(
		mbPerMinuteToGibPerHour(quality.preferredSizeMbPerMinute),
		minimum,
		maximum
	);
	return { minimum, preferred, maximum };
}

export function nextSliderQuality(
	quality: QualitySizeSetting,
	field: SliderField,
	rawValue: string
) {
	const nextValue = clamp(Number.parseFloat(rawValue), 0, sliderMaxGibPerHour);
	if (!Number.isFinite(nextValue)) {
		return quality;
	}
	let { minimum, preferred, maximum } = sliderValues(quality);
	if (field === 'minimum') {
		minimum = clamp(nextValue, 0, Math.max(0, maximum - sliderHandleGap * 2));
		preferred = Math.max(preferred, minimum);
	} else if (field === 'preferred') {
		preferred = clamp(nextValue, minimum, maximum);
	} else {
		maximum = clamp(nextValue, minimum + sliderHandleGap * 2, sliderMaxGibPerHour);
		preferred = Math.min(preferred, maximum);
	}
	return {
		...quality,
		minimumSizeMbPerMinute: gibPerHourToMbPerMinute(minimum),
		preferredSizeMbPerMinute:
			preferred >= sliderMaxGibPerHour ? null : gibPerHourToMbPerMinute(preferred),
		maximumSizeMbPerMinute: maximum >= sliderMaxGibPerHour ? null : gibPerHourToMbPerMinute(maximum)
	};
}

export function gibValue(value: number) {
	return value.toFixed(2);
}

export function mbPerMinuteTitle(label: string, value: number) {
	return `${label}: ${gibPerHourToMbPerMinute(value).toFixed(2)} MB/m`;
}

export function labelOffset(value: number) {
	return `${(value / sliderMaxGibPerHour) * 100}%`;
}

export function activeTrackStyle(values: { minimum: number; maximum: number }) {
	return `left: ${labelOffset(values.minimum)}; width: ${(Math.max(values.maximum - values.minimum, 0) / sliderMaxGibPerHour) * 100}%`;
}

export function clamp(value: number, minimum: number, maximum: number) {
	return Math.min(Math.max(value, minimum), maximum);
}
