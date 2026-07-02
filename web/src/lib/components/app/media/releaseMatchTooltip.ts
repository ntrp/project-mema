import type { MatchInfo } from './releaseCandidateDisplay';
import type { MediaType } from '$lib/settings/types';

export interface TooltipField {
	label: string;
	value: string;
}

export interface TooltipSection {
	label: string;
	fields: TooltipField[];
}

export function parsedTooltipSections(info: MatchInfo, mediaType: MediaType): TooltipSection[] {
	const parsed = info.parsed;
	return [
		{
			label: 'Parsed release',
			fields: compactFields([
				mediaType === 'movie' ? field('Movie', parsed.release.movieTitle) : undefined,
				mediaType === 'series' ? field('Series', parsed.release.seriesTitle) : undefined,
				field('Year', parsed.release.year),
				field('Season', parsed.release.seasonNumber),
				field('Episode', parsed.release.episodeNumber),
				mediaType === 'series' ? field('Season pack', parsed.release.seasonPack) : undefined,
				field('Edition', parsed.release.edition),
				field('Release group', parsed.release.releaseGroup),
				field('Release hash', parsed.release.releaseHash)
			])
		},
		{
			label: 'Parsed quality',
			fields: compactFields([
				field('Quality', parsed.quality.quality),
				field('Source', parsed.quality.source),
				field('Resolution', parsed.quality.resolution),
				field('Video codec', parsed.quality.videoCodec),
				field('Audio codec', parsed.quality.audioCodec),
				field('Audio channels', parsed.quality.audioChannels),
				field('Version', parsed.quality.version),
				field('Proper', parsed.quality.proper),
				field('Repack', parsed.quality.repack),
				field('Real', parsed.quality.real)
			])
		},
		{
			label: 'Parsed details',
			fields: compactFields([
				field('Languages', parsed.languages.join(', ')),
				field('Release type', parsed.details.releaseType)
			])
		}
	].filter((section) => section.fields.length > 0);
}

function compactFields(fields: (TooltipField | undefined)[]) {
	return fields.filter((field): field is TooltipField => field !== undefined);
}

function field(label: string, value: string | number | boolean | null | undefined) {
	if (value === undefined || value === null || value === '') return undefined;
	return { label, value: typeof value === 'boolean' ? (value ? 'Yes' : 'No') : String(value) };
}
