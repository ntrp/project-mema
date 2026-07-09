import type { MediaFulfillmentActionRequest } from '$lib/settings/types';

export function mediaFulfillmentActionKey(request: MediaFulfillmentActionRequest) {
	const scope =
		request.trackId ??
		request.otherFileId ??
		request.externalSubtitleId ??
		request.filePath ??
		'global';
	return [
		request.operation,
		scope,
		request.targetType ?? '',
		(request.languageId ?? '').toLowerCase()
	].join('|');
}
