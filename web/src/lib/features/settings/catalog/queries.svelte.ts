import { createQuery } from '@tanstack/svelte-query';
import {
	listDownloadClients,
	listIndexers,
	listLibraryFolders,
	listLanguages,
	listMediaProfiles,
	listMetadataProviders,
	listPathMappings,
	listSubtitleProviders,
	listCustomFormats,
	listTags,
	listUsers
} from './api';

export const settingsCatalogKeys = {
	all: ['settings-catalog'] as const,
	languages: () => [...settingsCatalogKeys.all, 'languages'] as const,
	tags: () => [...settingsCatalogKeys.all, 'tags'] as const,
	users: () => [...settingsCatalogKeys.all, 'users'] as const,
	downloadClients: () => [...settingsCatalogKeys.all, 'download-clients'] as const,
	indexers: () => [...settingsCatalogKeys.all, 'indexers'] as const,
	metadataProviders: () => [...settingsCatalogKeys.all, 'metadata-providers'] as const,
	subtitleProviders: () => [...settingsCatalogKeys.all, 'subtitle-providers'] as const,
	libraryFolders: () => [...settingsCatalogKeys.all, 'library-folders'] as const,
	pathMappings: () => [...settingsCatalogKeys.all, 'path-mappings'] as const,
	mediaProfiles: () => [...settingsCatalogKeys.all, 'media-profiles'] as const,
	customFormats: () => [...settingsCatalogKeys.all, 'custom-formats'] as const
};

export const createLanguagesQuery = (enabled: () => boolean = () => true) =>
	createQuery(() => ({
		queryKey: settingsCatalogKeys.languages(),
		queryFn: listLanguages,
		enabled: enabled()
	}));
export const createTagsQuery = (enabled: () => boolean = () => true) =>
	createQuery(() => ({
		queryKey: settingsCatalogKeys.tags(),
		queryFn: listTags,
		enabled: enabled()
	}));
export const createUsersQuery = (enabled: () => boolean = () => true) =>
	createQuery(() => ({
		queryKey: settingsCatalogKeys.users(),
		queryFn: listUsers,
		enabled: enabled()
	}));
export const createDownloadClientsQuery = (enabled: () => boolean) =>
	createQuery(() => ({
		queryKey: settingsCatalogKeys.downloadClients(),
		queryFn: listDownloadClients,
		enabled: enabled()
	}));
export const createIndexersQuery = (enabled: () => boolean) =>
	createQuery(() => ({
		queryKey: settingsCatalogKeys.indexers(),
		queryFn: listIndexers,
		enabled: enabled()
	}));
export const createMetadataProvidersQuery = (enabled: () => boolean) =>
	createQuery(() => ({
		queryKey: settingsCatalogKeys.metadataProviders(),
		queryFn: listMetadataProviders,
		enabled: enabled()
	}));
export const createSubtitleProvidersQuery = (enabled: () => boolean) =>
	createQuery(() => ({
		queryKey: settingsCatalogKeys.subtitleProviders(),
		queryFn: listSubtitleProviders,
		enabled: enabled()
	}));
export const createLibraryFoldersQuery = (enabled: () => boolean) =>
	createQuery(() => ({
		queryKey: settingsCatalogKeys.libraryFolders(),
		queryFn: listLibraryFolders,
		enabled: enabled()
	}));
export const createPathMappingsQuery = (enabled: () => boolean) =>
	createQuery(() => ({
		queryKey: settingsCatalogKeys.pathMappings(),
		queryFn: listPathMappings,
		enabled: enabled()
	}));
export const createMediaProfilesQuery = (enabled: () => boolean) =>
	createQuery(() => ({
		queryKey: settingsCatalogKeys.mediaProfiles(),
		queryFn: listMediaProfiles,
		enabled: enabled()
	}));
export const createCustomFormatsQuery = (enabled: () => boolean) =>
	createQuery(() => ({
		queryKey: settingsCatalogKeys.customFormats(),
		queryFn: listCustomFormats,
		enabled: enabled()
	}));
