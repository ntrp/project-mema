import type {
	LibraryMediaKind,
	LibraryScan,
	LibraryScanImportRequest,
	MediaSearchResult,
	MetadataProvider,
	QualityProfileOption
} from '$lib/settings/types';

export interface LibraryScanImportTableProps {
	scan: LibraryScan;
	qualityProfiles: QualityProfileOption[];
	metadataProviders: MetadataProvider[];
	loading: boolean;
	onSearchMatch: (
		_kind: LibraryMediaKind,
		_query: string,
		_providerId?: string
	) => Promise<MediaSearchResult[]>;
	onImport: (_scan: LibraryScan, _request: LibraryScanImportRequest) => Promise<void>;
}
