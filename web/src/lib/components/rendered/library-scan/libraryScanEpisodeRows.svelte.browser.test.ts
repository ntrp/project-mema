import * as browser from 'vitest/browser';
import { mount, tick, unmount } from 'svelte';
import { afterEach, describe, expect, it, vi } from 'vitest';

import RenderWithTooltip from '$lib/components/rendered/RenderWithTooltip.svelte';
import LibraryScanImportTable from '$lib/components/settings/library/scan/LibraryScanImportTable.svelte';
import LibraryScanImportTableBody from '$lib/components/settings/library/scan/LibraryScanImportTableBody.svelte';
import type { MatchDraft } from '$lib/components/settings/library/scan/libraryScanImport';
import type { LibraryScan, LibraryScanItem } from '$lib/settings/types';

const { page } = browser as unknown as typeof import('@vitest/browser/context');
let component: ReturnType<typeof mount> | undefined;
const target = document.createElement('div');

afterEach(async () => {
	if (component) await unmount(component);
	target.replaceChildren();
	component = undefined;
});

function episode(id: string, episodeNumber: number): LibraryScanItem {
	return {
		id,
		path: `/library/Show/Season 1/0${episodeNumber}.mkv`,
		fileName: `0${episodeNumber}.mkv`,
		seriesFolderPath: '/library/Show',
		detectedTitle: 'Show',
		detectedMediaKind: 'series',
		seasonNumber: 1,
		episodeNumber,
		status: 'pending',
		imported: false
	} as LibraryScanItem;
}

function draft(item: LibraryScanItem, selected = false): MatchDraft {
	return {
		selected,
		query: 'Show',
		mediaKind: 'series',
		metadataProviderId: '',
		matched: { title: 'Show', type: 'serie' },
		results: [],
		searching: false,
		searched: true,
		qualityProfileId: 'profile',
		monitorMode: 'all_episodes',
		minimumAvailability: 'released',
		seriesType: 'standard',
		seasonNumber: item.seasonNumber,
		episodeNumber: item.episodeNumber,
		removeDuplicate: false
	};
}

function renderRows(selected = false) {
	const rows = [episode('one', 1), episode('two', 2), episode('three', 3)];
	const drafts = $state(Object.fromEntries(rows.map((row) => [row.id, draft(row, selected)])));
	document.body.append(target);
	component = mount(LibraryScanImportTableBody, {
		target,
		props: {
			rows,
			folderPath: '/library',
			drafts,
			duplicateStates: {},
			qualityProfiles: [],
			metadataProviders: [],
			importingItemId: '',
			resettingItemId: '',
			onSearch: vi.fn(),
			onSelect: vi.fn(),
			onProviderChange: vi.fn(),
			onSeasonChange: vi.fn(),
			onEpisodeChange: vi.fn(),
			onResetImport: vi.fn()
		}
	});
	return { rows, drafts };
}

function pointer(element: EventTarget, type: string) {
	element.dispatchEvent(new PointerEvent(type, { bubbles: true, button: 0, pointerId: 1 }));
}

describe('series episode selection', () => {
	it('selects a series from its root checkbox', async () => {
		const { drafts } = renderRows();
		await tick();
		const checkboxes = target.querySelectorAll('[role="checkbox"]');
		expect(checkboxes).toHaveLength(4);
		(checkboxes[0] as HTMLElement).click();
		await tick();
		expect(Object.values(drafts).map((row) => row.selected)).toEqual([true, true, true]);
	});

	it('collapses and expands episode files from the series row', async () => {
		renderRows();
		await tick();
		const toggle = target.querySelector('button[aria-expanded]');
		expect(target.textContent).toContain('Season 1/01.mkv');
		toggle?.closest('tr')?.dispatchEvent(new MouseEvent('click', { bubbles: true }));
		await tick();
		expect(toggle?.getAttribute('aria-expanded')).toBe('false');
		expect(target.textContent).not.toContain('Season 1/01.mkv');
		toggle?.dispatchEvent(new MouseEvent('click', { bubbles: true }));
		await tick();
		expect(target.textContent).toContain('Season 1/01.mkv');
	});

	it('disables episode selection until season and episode are set', async () => {
		const { drafts } = renderRows();
		drafts.one.seasonNumber = undefined;
		await tick();
		const checkbox = target.querySelectorAll('[role="checkbox"]')[1];
		expect(checkbox.getAttribute('data-disabled')).not.toBeNull();
		drafts.one.seasonNumber = 1;
		await tick();
		expect(checkbox.getAttribute('data-disabled')).toBeNull();
	});

	it('toggles an episode once for a pointer click', async () => {
		const { drafts } = renderRows();
		await tick();
		const checkbox = target.querySelectorAll('[role="checkbox"]')[1];
		pointer(checkbox, 'pointerdown');
		pointer(checkbox, 'pointerup');
		(checkbox as HTMLElement).click();
		await tick();
		expect(drafts.one.selected).toBe(true);
	});

	it('changes season and episode values only on the edited row', async () => {
		const rows = [episode('one', 1), episode('two', 2)];
		for (const row of rows) {
			row.matchedTitle = 'Show';
			row.matchedMediaKind = 'series';
		}
		document.body.append(target);
		component = mount(RenderWithTooltip, {
			target,
			props: {
				component: LibraryScanImportTable as never,
				componentProps: {
					scan: {
						id: 'scan',
						folderId: 'folder',
						folderPath: '/library',
						folderKind: 'series',
						totalFiles: rows.length,
						items: rows
					} as LibraryScan,
					qualityProfiles: [{ id: 'profile', name: 'Default' }],
					metadataProviders: [],
					loading: false,
					onSearchMatch: vi.fn().mockResolvedValue([]),
					onImport: vi.fn(),
					onResetImport: vi.fn()
				}
			}
		});
		await tick();
		await page.getByText('Import Selected').hover();
		await vi.waitFor(() =>
			expect(document.body.textContent).toContain('Select at least one file to import')
		);
		const checkboxes = target.querySelectorAll('[aria-label^="Select Season"]');
		expect([...checkboxes].map((checkbox) => checkbox.getAttribute('aria-checked'))).toEqual([
			'false',
			'false'
		]);
		(checkboxes[1] as HTMLElement).click();
		await tick();
		await page.getByRole('button', { name: 'Select profile' }).first().click();
		await page.getByRole('option', { name: 'Default' }).click();
		await tick();
		const importButton = [...target.querySelectorAll('button')].find((button) =>
			button.textContent?.includes('Import Selected')
		);
		expect(importButton?.disabled).toBe(false);
		(checkboxes[0] as HTMLElement).click();
		await tick();
		const inputs = target.querySelectorAll<HTMLInputElement>('input[type="number"]');
		inputs[0].value = '2';
		inputs[0].dispatchEvent(new Event('change', { bubbles: true }));
		await tick();
		expect([...inputs].map((input) => input.value)).toEqual(['2', '', '1', '2']);
	});

	it('selects and unselects episodes by dragging across checkboxes', async () => {
		let setup = renderRows();
		await tick();
		let checkboxes = target.querySelectorAll('[role="checkbox"]');
		pointer(checkboxes[1], 'pointerdown');
		pointer(checkboxes[2], 'pointerenter');
		pointer(window, 'pointerup');
		expect(setup.drafts.one.selected).toBe(true);
		expect(setup.drafts.two.selected).toBe(true);
		expect(setup.drafts.three.selected).toBe(false);

		await unmount(component!);
		target.replaceChildren();
		component = undefined;
		setup = renderRows(true);
		await tick();
		checkboxes = target.querySelectorAll('[role="checkbox"]');
		pointer(checkboxes[1], 'pointerdown');
		pointer(checkboxes[2], 'pointerenter');
		pointer(window, 'pointerup');
		expect(setup.drafts.one.selected).toBe(false);
		expect(setup.drafts.two.selected).toBe(false);
		expect(setup.drafts.three.selected).toBe(true);
	});
});
