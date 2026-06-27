import type { components } from '$lib/api/generated/schema';

export type DownloadClient = components['schemas']['DownloadClient'];
export type DownloadClientRequest = components['schemas']['DownloadClientRequest'];
export type DownloadClientType = components['schemas']['DownloadClientType'];
export type Indexer = components['schemas']['Indexer'];
export type IndexerRequest = components['schemas']['IndexerRequest'];
export type IndexerType = components['schemas']['IndexerType'];
export type IntegrationTestResponse = components['schemas']['IntegrationTestResponse'];

export type DownloadClientForm = DownloadClientRequest & { id?: string };
export type IndexerForm = Omit<IndexerRequest, 'categories'> & {
	id?: string;
	categoriesText: string;
};

export type AppView = 'home' | 'settings';
export type HomeSection = 'explore' | 'movies' | 'series' | 'activity';
export type SettingsSection = 'download-clients' | 'indexers';
export interface SettingsData {
	downloadClients: DownloadClient[];
	indexers: Indexer[];
}

export type IntegrationTestResults = Record<string, IntegrationTestResponse | undefined>;
