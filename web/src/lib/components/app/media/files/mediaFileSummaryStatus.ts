import type { MediaFileRow } from '$lib/components/app/media/files/mediaFiles';

export type MediaFileSummaryStatus = NonNullable<MediaFileRow['requirements']>['audio'];

export function fallbackRequirementStatus(label: string, exists: boolean): MediaFileSummaryStatus {
	return {
		state: exists ? 'ignored' : 'missing',
		label: exists ? 'Ignored' : 'Missing',
		details: [`${label} state was not provided by the backend`]
	};
}
