type QualityLike = {
	qualityId: string;
	name: string;
};

export type QualityResolutionGroup<TQuality extends QualityLike> = {
	id: string;
	label: string;
	qualities: TQuality[];
};

const resolutionGroups = [
	{ id: 'other', label: 'Other' },
	{ id: 'sd', label: 'SD' },
	{ id: '480p', label: '480p' },
	{ id: '576p', label: '576p' },
	{ id: '720p', label: '720p' },
	{ id: '1080p', label: '1080p' },
	{ id: '4k', label: '4K' },
	{ id: 'native', label: 'Native' }
] as const;

export function groupQualitiesByResolution<TQuality extends QualityLike>(
	qualities: TQuality[]
): QualityResolutionGroup<TQuality>[] {
	const groups = resolutionGroups.map((group) => ({
		...group,
		qualities: [] as TQuality[]
	}));

	for (const quality of qualities) {
		const group = groups.find((item) => item.id === resolutionGroupId(quality)) ?? groups[0];
		group.qualities.push(quality);
	}

	return groups.filter((group) => group.qualities.length > 0);
}

function resolutionGroupId(quality: QualityLike) {
	const value = `${quality.qualityId} ${quality.name}`.toLowerCase();
	if (quality.qualityId === 'br-disk' || quality.qualityId === 'raw-hd') {
		return 'native';
	}
	if (value.includes('2160') || value.includes('4k')) {
		return '4k';
	}
	if (value.includes('1080')) {
		return '1080p';
	}
	if (value.includes('720')) {
		return '720p';
	}
	if (value.includes('576')) {
		return '576p';
	}
	if (value.includes('480')) {
		return '480p';
	}
	if (/cam|dvd|regional|sdtv|telecine|telesync/.test(value)) {
		return 'sd';
	}
	return 'other';
}
