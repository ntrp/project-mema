import type {
	Language,
	MediaItem,
	ReleaseCandidate,
	ReleaseOverrideDetails
} from '$lib/settings/types';
import type { ReleaseSearchContext } from './releaseSearchQuery';

export interface MediaFileSearchModalProps {
	item: MediaItem;
	grabbingKey?: string;
	searchContext?: ReleaseSearchContext;
	languages: Language[];
	canManage: boolean;
	onGrab: (
		_item: MediaItem,
		_release: ReleaseCandidate,
		_overrideMatch?: boolean,
		_details?: ReleaseOverrideDetails
	) => void;
	onClose: () => void;
}
