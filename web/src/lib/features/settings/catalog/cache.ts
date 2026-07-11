import type { QueryClient } from '@tanstack/svelte-query';
import type {
	DownloadClient,
	CustomFormat,
	Indexer,
	Language,
	LibraryFolder,
	ManagedUser,
	MediaProfile,
	MetadataProvider,
	SubtitleProvider,
	PathMapping,
	Tag
} from '$lib/settings/types';
import { settingsCatalogKeys } from './queries.svelte';

export function createSettingsCatalogCache(client: QueryClient) {
	const removeLanguage = (code: string) =>
		client.setQueryData<Language[]>(settingsCatalogKeys.languages(), (items) =>
			(items ?? []).filter((item) => item.code !== code)
		);
	const removeTag = (id: string) =>
		client.setQueryData<Tag[]>(settingsCatalogKeys.tags(), (items) =>
			(items ?? []).filter((item) => item.id !== id)
		);
	const removeUser = (id: string) =>
		client.setQueryData<ManagedUser[]>(settingsCatalogKeys.users(), (items) =>
			(items ?? []).filter((item) => item.id !== id)
		);
	const removeDownloadClient = (id: string) =>
		client.setQueryData<DownloadClient[]>(settingsCatalogKeys.downloadClients(), (items) =>
			(items ?? []).filter((item) => item.id !== id)
		);
	const removeIndexer = (id: string) =>
		client.setQueryData<Indexer[]>(settingsCatalogKeys.indexers(), (items) =>
			(items ?? []).filter((item) => item.id !== id)
		);
	const removeMetadataProvider = (id: string) =>
		client.setQueryData<MetadataProvider[]>(settingsCatalogKeys.metadataProviders(), (items) =>
			(items ?? []).filter((item) => item.id !== id)
		);
	const removeSubtitleProvider = (id: string) =>
		client.setQueryData<SubtitleProvider[]>(settingsCatalogKeys.subtitleProviders(), (items) =>
			(items ?? []).filter((item) => item.id !== id)
		);
	const update = <T extends { id: string }>(key: readonly unknown[], item: T) =>
		client.setQueryData<T[]>(key, (items) => [
			item,
			...(items ?? []).filter((entry) => entry.id !== item.id)
		]);
	const upsertLibraryFolder = (item: LibraryFolder) =>
		update(settingsCatalogKeys.libraryFolders(), item);
	const upsertPathMapping = (item: PathMapping) => update(settingsCatalogKeys.pathMappings(), item);
	const removeLibraryFolder = (id: string) =>
		client.setQueryData<LibraryFolder[]>(settingsCatalogKeys.libraryFolders(), (items) =>
			(items ?? []).filter((item) => item.id !== id)
		);
	const removePathMapping = (id: string) =>
		client.setQueryData<PathMapping[]>(settingsCatalogKeys.pathMappings(), (items) =>
			(items ?? []).filter((item) => item.id !== id)
		);
	const removeMediaProfile = (id: string) =>
		client.setQueryData<MediaProfile[]>(settingsCatalogKeys.mediaProfiles(), (items) =>
			(items ?? []).filter((item) => item.id !== id)
		);
	const removeCustomFormat = (id: string) =>
		client.setQueryData<CustomFormat[]>(settingsCatalogKeys.customFormats(), (items) =>
			(items ?? []).filter((item) => item.id !== id)
		);
	const refresh = () => client.invalidateQueries({ queryKey: settingsCatalogKeys.all });
	const clear = () => client.removeQueries({ queryKey: settingsCatalogKeys.all });
	return {
		removeLanguage,
		removeTag,
		removeUser,
		removeDownloadClient,
		removeIndexer,
		removeMetadataProvider,
		removeSubtitleProvider,
		upsertLibraryFolder,
		upsertPathMapping,
		removeLibraryFolder,
		removePathMapping,
		removeMediaProfile,
		removeCustomFormat,
		refresh,
		clear
	};
}
