import { mount, tick, unmount } from 'svelte';
import { afterEach, describe, expect, it, vi } from 'vitest';

import MetadataProviderCard from './MetadataProviderCard.svelte';
import type { MetadataProviderForm } from '$lib/settings/types';

const anidbDefinition = {
	type: 'anidb' as const,
	name: 'AniDB',
	baseUrl: 'http://api.anidb.net:9001',
	priority: 120,
	fields: 'anidb' as const
};
const anilistDefinition = {
	type: 'anilist' as const,
	name: 'AniList',
	baseUrl: 'https://graphql.anilist.co',
	priority: 120,
	fields: 'anilist' as const
};
const tmdbDefinition = {
	type: 'tmdb' as const,
	name: 'TMDB',
	baseUrl: 'https://api.themoviedb.org/3',
	priority: 100,
	fields: 'tmdb' as const
};
let component: ReturnType<typeof mount> | undefined;
let target = document.createElement('div');

afterEach(async () => {
	if (component) await unmount(component);
	target.remove();
	component = undefined;
});

function input(selector: string) {
	const element = target.querySelector<HTMLInputElement>(selector);
	if (!element) throw new Error(`missing input: ${selector}`);
	return element;
}

function updateInput(change: { element: HTMLInputElement; value: string }) {
	change.element.value = change.value;
	change.element.dispatchEvent(new window.InputEvent('input', { bubbles: true }));
}

function click(selector: string) {
	const element = target.querySelector(selector);
	if (!(element instanceof window.HTMLElement)) throw new Error(`missing element: ${selector}`);
	element.click();
}

describe('MetadataProviderCard interactions', () => {
	it('enables TMDB fallback before save', async () => {
		const onSave = vi.fn<(form: MetadataProviderForm) => void>();
		target = document.createElement('div');
		document.body.append(target);
		component = mount(MetadataProviderCard, {
			target,
			props: { definition: tmdbDefinition, onSave, onTest: vi.fn() }
		});
		await tick();

		const checkboxes = target.querySelectorAll('[role="checkbox"]');
		expect(checkboxes).toHaveLength(2);
		(checkboxes.item(1) as HTMLElement).click();
		click('button[type="submit"]');

		expect(onSave).toHaveBeenCalledWith(expect.objectContaining({ fallbackEnabled: true }));
	});

	it('updates AniDB inputs, enabled state, and priority before save', async () => {
		const onSave = vi.fn<(form: MetadataProviderForm) => void>();
		target = document.createElement('div');
		document.body.append(target);
		component = mount(MetadataProviderCard, {
			target,
			props: { definition: anidbDefinition, onSave, onTest: vi.fn() }
		});
		await tick();

		updateInput({ element: input('input[maxlength="2000"]'), value: 'http://metadata.test/anidb' });
		updateInput({ element: input('input[pattern="[a-z0-9]+"]'), value: 'registeredclient1' });
		updateInput({ element: input('input[type="number"][required]'), value: '3' });
		updateInput({ element: input('input[type="number"][max="1000"]'), value: '321' });
		click('[role="checkbox"]');
		await tick();
		expect(target.querySelector('[role="checkbox"]')?.getAttribute('aria-checked')).toBe('true');
		const form = target.querySelector('form');
		expect(form?.checkValidity()).toBe(true);
		form?.requestSubmit();

		expect(onSave).toHaveBeenCalledWith(
			expect.objectContaining({
				baseUrl: 'http://metadata.test/anidb',
				apiKey: 'registeredclient1',
				pin: '3',
				enabled: true,
				priority: 321
			})
		);
	});

	it('handles configured enabled, empty priority, save, and test actions', async () => {
		const onSave = vi.fn<(form: MetadataProviderForm) => void>();
		const onTest = vi.fn();
		target = document.createElement('div');
		document.body.append(target);
		component = mount(MetadataProviderCard, {
			target,
			props: {
				definition: anidbDefinition,
				provider: {
					id: 'provider-anidb',
					name: 'AniDB',
					type: 'anidb',
					baseUrl: anidbDefinition.baseUrl,
					apiKey: 'registeredclient1',
					pin: '2',
					apiKeySet: true,
					pinSet: true,
					accessTokenSet: false,
					enabled: true,
					fallbackEnabled: false,
					priority: 120,
					createdAt: '2026-07-03T00:00:00Z',
					updatedAt: '2026-07-03T00:00:00Z'
				},
				onSave,
				onTest
			}
		});
		await tick();

		click('[role="checkbox"]');
		updateInput({ element: input('input[type="number"][max="1000"]'), value: '' });
		await tick();
		const buttons = target.querySelectorAll('button');
		expect(Array.from(buttons, (button) => button.textContent?.trim())).toEqual([
			'',
			'Save',
			'Test'
		]);
		buttons.item(buttons.length - 1).click();
		click('button[type="submit"]');

		expect(onTest).toHaveBeenCalledWith('provider-anidb');
		expect(onSave).toHaveBeenCalledWith(expect.objectContaining({ enabled: false, priority: 0 }));
	});

	it('toggles, prioritizes, saves, and tests AniList without credentials', async () => {
		const onSave = vi.fn<(form: MetadataProviderForm) => void>();
		const onTest = vi.fn();
		target = document.createElement('div');
		document.body.append(target);
		component = mount(MetadataProviderCard, {
			target,
			props: {
				definition: anilistDefinition,
				provider: {
					id: 'provider-anilist',
					name: 'AniList',
					type: 'anilist',
					baseUrl: anilistDefinition.baseUrl,
					apiKeySet: false,
					pinSet: false,
					accessTokenSet: false,
					enabled: true,
					fallbackEnabled: false,
					priority: 120,
					createdAt: '2026-07-24T00:00:00Z',
					updatedAt: '2026-07-24T00:00:00Z'
				},
				onSave,
				onTest
			}
		});
		await tick();

		expect(target.textContent).toContain('No credentials required');
		expect(target.querySelectorAll('input')).toHaveLength(2);
		click('[role="checkbox"]');
		updateInput({ element: input('input[type="number"][max="1000"]'), value: '125' });
		click('button[type="button"]:last-child');
		click('button[type="submit"]');

		expect(onTest).toHaveBeenCalledWith('provider-anilist');
		expect(onSave).toHaveBeenCalledWith(
			expect.objectContaining({ type: 'anilist', enabled: false, priority: 125 })
		);
		const saved = onSave.mock.calls[0][0];
		expect(saved).not.toHaveProperty('apiKey');
		expect(saved).not.toHaveProperty('pin');
		expect(saved).not.toHaveProperty('accessToken');
	});
});
