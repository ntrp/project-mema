import { mount, tick, unmount } from 'svelte';
import { afterEach, describe, expect, it, vi } from 'vitest';

import MetadataProviderCard from './MetadataProviderCard.svelte';
import type { MetadataProviderForm } from '$lib/settings/types';

const definition = {
	type: 'anidb' as const,
	name: 'AniDB',
	baseUrl: 'http://api.anidb.net:9001',
	priority: 120,
	fields: 'anidb' as const
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

function update(element: HTMLInputElement, value: string) {
	element.value = value;
	element.dispatchEvent(new window.InputEvent('input', { bubbles: true }));
}

function click(selector: string) {
	const element = target.querySelector(selector);
	if (!(element instanceof window.HTMLElement)) throw new Error(`missing element: ${selector}`);
	element.click();
}

describe('MetadataProviderCard interactions', () => {
	it('updates AniDB inputs, enabled state, and priority before save', async () => {
		const onSave = vi.fn<(form: MetadataProviderForm) => void>();
		target = document.createElement('div');
		document.body.append(target);
		component = mount(MetadataProviderCard, {
			target,
			props: { definition, onSave, onTest: vi.fn() }
		});
		await tick();

		update(input('input[maxlength="2000"]'), 'http://metadata.test/anidb');
		update(input('input[pattern="[a-z0-9]+"]'), 'registeredclient1');
		update(input('input[type="number"][required]'), '3');
		update(input('input[type="number"][max="1000"]'), '321');
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
				definition,
				provider: {
					id: 'provider-anidb',
					name: 'AniDB',
					type: 'anidb',
					baseUrl: definition.baseUrl,
					apiKey: 'registeredclient1',
					pin: '2',
					apiKeySet: true,
					pinSet: true,
					accessTokenSet: false,
					enabled: true,
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
		update(input('input[type="number"][max="1000"]'), '');
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
});
