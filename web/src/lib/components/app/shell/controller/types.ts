import type { AppView, HomeSection, SettingsSection, SystemSection } from '$lib/settings/types';

export interface ServerEventEnvelope<T = unknown> {
	id: string;
	type: string;
	time: string;
	data: T;
}

export interface AppShellOptions {
	initialView?: AppView;
	initialHomeSection?: HomeSection;
	initialSettingsSection?: SettingsSection;
	initialSystemSection?: SystemSection;
	initialSelectedMediaItemId?: string;
	initialSelectedRequestId?: string;
	initialAdvancedQuery?: string;
	initialMetadataProvider?: string;
	initialMetadataType?: string;
	initialMetadataExternalId?: string;
	initialCollectionProvider?: string;
	initialCollectionId?: string;
	initialDiscoverSectionId?: string;
	initialRelatedSectionKind?: RelatedSectionKind;
	initialPeopleSectionKind?: PeopleSectionKind;
}

export type RelatedSectionKind = 'recommendations' | 'similar';
export type PeopleSectionKind = 'cast' | 'crew';
