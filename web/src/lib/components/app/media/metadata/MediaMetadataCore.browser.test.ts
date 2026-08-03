import { mount, unmount } from 'svelte';
import { afterEach, describe, expect, it } from 'vitest';
import type { MediaMetadataDetails } from '$lib/settings/types';
import MediaMetadataCore from './MediaMetadataCore.svelte';

const browserDescribe = typeof document === 'undefined' ? describe.skip : describe;

browserDescribe('MediaMetadataCore', () => {
	let component: object | undefined;

	afterEach(async () => {
		if (component) await unmount(component);
		document.body.replaceChildren();
	});

	it('renders duplicate provider cast entries once without duplicate keyed blocks', () => {
		component = mount(MediaMetadataCore, {
			target: document.body,
			props: {
				detail: {
					title: 'Shingeki no Kyojin',
					type: 'serie',
					overview: 'Test overview',
					cast: [
						{ name: 'Kaji Yuuki', role: 'Eren Yeager' },
						{ name: 'Kaji Yuuki', role: 'Eren Yeager' }
					]
				} as MediaMetadataDetails
			}
		});

		expect(
			[...document.querySelectorAll('strong')].filter(
				(element) => element.textContent === 'Kaji Yuuki'
			)
		).toHaveLength(1);
	});
});
