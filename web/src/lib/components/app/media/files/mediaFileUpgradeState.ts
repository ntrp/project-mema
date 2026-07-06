import type { MediaFileProfileOption } from '$lib/components/app/media/file-data/mediaFileProfiles';

export type MediaFileUpgradeState = 'current' | 'upgradeable' | 'blocked' | 'missing';

export interface MediaFileUpgradeInfo {
	state: MediaFileUpgradeState;
	label: string;
	reasons: string[];
}

export function mediaFileUpgradeInfo(
	exists: boolean,
	quality: string,
	formats: string[],
	profile?: MediaFileProfileOption
): MediaFileUpgradeInfo {
	if (!exists) return { state: 'missing', label: 'Missing', reasons: ['File is missing'] };
	const qualityIds = profile?.qualityIds ?? [];
	if (!profile || qualityIds.length === 0) {
		return { state: 'current', label: 'Current', reasons: ['No quality target configured'] };
	}
	if (profile.upgradesAllowed === false) {
		return { state: 'blocked', label: 'Blocked', reasons: ['Upgrades are disabled'] };
	}
	const current = qualityID(quality, formats);
	if (!current) {
		return { state: 'blocked', label: 'Blocked', reasons: ['File quality could not be detected'] };
	}
	const currentIndex = qualityIds.indexOf(current);
	if (currentIndex < 0) {
		return {
			state: 'blocked',
			label: 'Blocked',
			reasons: [`${current} is not enabled in profile`]
		};
	}
	const target = profile.upgradeUntilQualityId;
	const targetIndex = target ? qualityIds.indexOf(target) : -1;
	if (target && targetIndex > currentIndex) {
		return { state: 'upgradeable', label: 'Upgradeable', reasons: [`Upgrade target is ${target}`] };
	}
	return { state: 'current', label: 'Current', reasons: ['At or above upgrade target'] };
}

function qualityID(quality: string, formats: string[]) {
	const normalizedQuality = quality.toLowerCase();
	const normalizedFormats = formats.map((format) => format.toLowerCase());
	const source =
		normalizedFormats.includes('web-dl') || normalizedQuality.startsWith('webdl-')
			? 'webdl'
			: normalizedFormats.includes('bluray') || normalizedQuality.startsWith('bluray-')
				? 'bluray'
				: '';
	const resolution = normalizedQuality.match(/\b(2160p|1080p|720p|576p|480p)\b/)?.[1];
	if (!source || !resolution) return undefined;
	return `${source}-${resolution}`;
}
