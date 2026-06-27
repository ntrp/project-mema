import type {
	DownloadClient,
	DownloadClientForm,
	DownloadClientRequest,
	Indexer,
	IndexerForm,
	IndexerRequest
} from './types';

export function emptyDownloadClientForm(): DownloadClientForm {
	return {
		name: '',
		type: 'transmission',
		baseUrl: '',
		username: '',
		password: '',
		apiKey: '',
		category: '',
		enabled: true,
		priority: 100
	};
}

export function emptyIndexerForm(): IndexerForm {
	return {
		name: '',
		type: 'torznab',
		baseUrl: '',
		apiKey: '',
		categoriesText: '',
		enabled: true,
		priority: 100
	};
}

export function downloadClientFormFromClient(client: DownloadClient): DownloadClientForm {
	return {
		id: client.id,
		name: client.name,
		type: client.type,
		baseUrl: client.baseUrl,
		username: client.username ?? '',
		password: client.password ?? '',
		apiKey: client.apiKey ?? '',
		category: client.category ?? '',
		enabled: client.enabled,
		priority: client.priority
	};
}

export function indexerFormFromIndexer(indexer: Indexer): IndexerForm {
	return {
		id: indexer.id,
		name: indexer.name,
		type: indexer.type,
		baseUrl: indexer.baseUrl,
		apiKey: indexer.apiKey ?? '',
		categoriesText: (indexer.categories ?? []).join(', '),
		enabled: indexer.enabled,
		priority: indexer.priority
	};
}

export function normalizeDownloadClientForm(form: DownloadClientForm): DownloadClientRequest {
	return {
		name: form.name.trim(),
		type: form.type,
		baseUrl: form.baseUrl.trim(),
		username: optionalString(form.username),
		password: optionalString(form.password),
		apiKey: optionalString(form.apiKey),
		category: optionalString(form.category),
		enabled: form.enabled,
		priority: form.priority
	};
}

export function normalizeIndexerForm(form: IndexerForm): IndexerRequest {
	return {
		name: form.name.trim(),
		type: form.type,
		baseUrl: form.baseUrl.trim(),
		apiKey: optionalString(form.apiKey),
		categories: parseCategories(form.categoriesText),
		enabled: form.enabled,
		priority: form.priority
	};
}

function optionalString(value: string | undefined) {
	const trimmed = value?.trim() ?? '';
	return trimmed === '' ? undefined : trimmed;
}

function parseCategories(value: string) {
	return value
		.split(',')
		.map((item) => Number.parseInt(item.trim(), 10))
		.filter((item) => Number.isInteger(item));
}
