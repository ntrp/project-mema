import { render } from 'svelte/server';
import { describe, expect, it, vi } from 'vitest';

import LibraryScanImportRow from '$lib/components/settings/library/scan/LibraryScanImportRow.svelte';
import LibraryScanImportTable from '$lib/components/settings/library/scan/LibraryScanImportTable.svelte';
import FileNamingSettings from '$lib/components/settings/library/FileNamingSettings.svelte';
import QualitySizeRow from '$lib/components/settings/quality/QualitySizeRow.svelte';
import SystemStatusSettings from '$lib/components/settings/system/SystemStatusSettings.svelte';
import SystemLogFilesSettings from '$lib/components/settings/system/logs/SystemLogFilesSettings.svelte';
import type {
	LibraryScan,
	LibraryScanItem,
	MediaSearchResult,
	MetadataProvider,
	QualityProfileOption,
	QualitySizeSetting
} from '$lib/settings/types';
import type { MatchDraft } from '$lib/components/settings/library/scan/libraryScanImport';
import { renderWithTooltip } from './renderHelpers';

describe('rendered library settings (SCN-LIBRARY-004)', () => {
	it('renders file naming defaults and examples', () => {
		const { body } = renderWithTooltip(FileNamingSettings, {});

		expect(body).toContain('File Naming');
		expect(body).toContain('Movie');
		expect(body).toContain('Series');
		expect(body).toContain('Main folder');
		expect(body).toContain('Save templates');
		expect(body).toContain('Defaults');
	});

	it('renders scan import controls and matched row choices', () => {
		const item = scanItem({ id: 'scan-item-1', fileName: 'Scenario.Movie.2026.mkv' });
		const profiles = [{ id: 'profile-1', name: 'Scenario Profile' }] as QualityProfileOption[];
		const providers = [
			{ id: 'metadata-1', name: 'TMDB', type: 'tmdb', enabled: true }
		] as MetadataProvider[];
		const draft = {
			selected: true,
			query: 'Scenario Movie',
			mediaKind: 'movie',
			metadataProviderId: 'metadata-1',
			matched: { title: 'Scenario Movie', type: 'movie', year: 2026 } as MediaSearchResult,
			results: [],
			searching: false,
			searched: true,
			qualityProfileId: 'profile-1',
			monitorMode: 'only_media',
			minimumAvailability: 'released',
			seriesType: 'standard',
			removeDuplicate: false
		} as MatchDraft;

		const table = render(LibraryScanImportTable, {
			props: {
				scan: {
					id: 'scan-1',
					folderPath: '/downloads',
					status: 'completed',
					totalFiles: 1,
					items: [item]
				} as LibraryScan,
				qualityProfiles: profiles,
				metadataProviders: providers,
				loading: false,
				onSearchMatch: vi.fn(),
				onImport: vi.fn(),
				onResetImport: vi.fn()
			}
		});
		const row = render(LibraryScanImportRow, {
			props: {
				item,
				folderPath: '/downloads',
				draft,
				qualityProfiles: profiles,
				metadataProviders: providers,
				onSearch: vi.fn(),
				onSelect: vi.fn(),
				onProviderChange: vi.fn(),
				onResetImport: vi.fn()
			}
		});
		const unmatchedRow = render(LibraryScanImportRow, {
			props: {
				item,
				folderPath: '/downloads',
				draft: { ...draft, matched: undefined, selected: false, searched: false },
				qualityProfiles: profiles,
				metadataProviders: providers,
				onSearch: vi.fn(),
				onSelect: vi.fn(),
				onProviderChange: vi.fn(),
				onResetImport: vi.fn()
			}
		});

		expect(table.body).toContain('1 files');
		expect(table.body).toContain('Directory / File');
		expect(table.body).toContain('Metadata provider');
		expect(table.body).toContain('Import Selected');
		expect(row.body).toContain('Scenario.Movie.2026/Scenario.Movie.2026.mkv');
		expect(row.body).not.toContain('/downloads');
		expect(row.body).toContain('Scenario Movie (2026)');
		expect(unmatchedRow.body).toContain('No match');
		expect(unmatchedRow.body).toContain('text-amber-500');
		expect(row.body).toContain('Scenario Profile');
		expect(row.body).toContain('Only this media');
	});
});

describe('rendered system and quality settings (SCN-SETTINGS-009)', () => {
	it('renders quality size fields and validation errors', () => {
		const { body } = renderWithTooltip(QualitySizeRow, {
			quality: {
				qualityId: 'q-1080p',
				name: 'HD-1080p',
				minimumSizeMbPerMinute: 10,
				preferredSizeMbPerMinute: 5,
				maximumSizeMbPerMinute: 20
			} as QualitySizeSetting,
			onChange: vi.fn()
		});

		expect(body).toContain('HD-1080p');
		expect(body).toContain('q-1080p');
		expect(body).toContain('Preferred must be at least minimum');
		expect(body).toContain('HD-1080p minimum size GiB per hour');
	});

	it('renders initial system status and log file states', () => {
		const status = renderWithTooltip(SystemStatusSettings, {});
		const logFiles = renderWithTooltip(SystemLogFilesSettings, {});

		expect(status.body).toContain('About');
		expect(status.body).toContain('Loading system status');
		expect(status.body).toContain('Refresh');
		expect(logFiles.body).toContain('Log files');
		expect(logFiles.body).toContain('Name');
		expect(logFiles.body).toContain('No log files retained.');
	});
});

function scanItem(overrides: Partial<LibraryScanItem>): LibraryScanItem {
	return {
		id: 'item-1',
		scanId: 'scan-1',
		path: '/downloads/Scenario.Movie.2026/Scenario.Movie.2026.mkv',
		fileName: 'Scenario.Movie.2026.mkv',
		status: 'pending',
		detectedTitle: 'Scenario Movie',
		detectedYear: 2026,
		detectedMediaKind: 'movie',
		...overrides
	} as LibraryScanItem;
}
