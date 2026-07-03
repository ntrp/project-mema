import { seasonNumberFromName } from '$lib/components/app/media/files/mediaFiles';
import type { MediaMetadataDetails, MediaMetadataSeason } from '$lib/settings/types';

export interface SeasonOption {
	value: string;
	label: string;
	season: MediaMetadataSeason;
}

export function seasonOptions(details?: MediaMetadataDetails): SeasonOption[] {
	return (details?.seasons ?? [])
		.map((season, index) => {
			const number = seasonNumberFromName(season.name) ?? index + 1;
			return { value: String(number), label: season.name, season };
		})
		.sort((left, right) => Number(left.value) - Number(right.value));
}

export function selectedSeason(options: SeasonOption[], seasonNumber: string) {
	return options.find((option) => option.value === seasonNumber);
}

export function episodeLabel(episode: { episodeNumber: number; name: string }) {
	return `E${String(episode.episodeNumber).padStart(2, '0')} ${episode.name}`;
}

export function episodeValueFromNumbers(values: number[]) {
	return Array.from(new Set(values))
		.filter((value) => Number.isInteger(value) && value > 0)
		.sort((left, right) => left - right)
		.join(', ');
}

export function episodeNumbers(value: string) {
	return value
		.split(/[\s,]+/)
		.filter((part) => part.trim() !== '')
		.map((part) => Number(part))
		.filter((part) => Number.isInteger(part) && part > 0);
}
