export interface MediaReleaseDates {
	releaseDate?: string;
	firstAirDate?: string;
}

export function isUnreleasedMedia(item: MediaReleaseDates) {
	const date = item.releaseDate ?? item.firstAirDate;
	if (!date) return false;
	const value = Date.parse(`${date}T00:00:00Z`);
	return Number.isFinite(value) && value > Date.now();
}
