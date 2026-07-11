import type { MediaItem, MediaRequest } from '$lib/settings/types';
import type { RunCommandMutation } from '$lib/app/query/commandMutation.svelte';

export interface MediaDeps {
	clearNotice: () => void;
	runMutation?: RunCommandMutation;
	loadMediaItems: () => Promise<void>;
	removeActivityForMedia: (_mediaItemId: string) => void;
	removeReleaseResults: (_mediaItemId: string) => void;
	mediaItems: () => MediaItem[];
	upsertMediaItem: (_item: MediaItem) => void;
	mapMediaItems: (_map: (_item: MediaItem) => MediaItem) => void;
	removeMediaItem: (_id: string) => void;
	upsertMediaRequest: (_request: MediaRequest) => void;
	mapMediaRequests: (_map: (_request: MediaRequest) => MediaRequest) => void;
}
