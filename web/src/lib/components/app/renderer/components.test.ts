import { describe, expect, it } from 'vitest';
import { render } from 'svelte/server';
import { renderWithTooltip } from '$lib/components/rendered/renderHelpers';
import PlayOnDevice from './PlayOnDevice.svelte';
import RendererRemote from './RendererRemote.svelte';

describe('play on device components', () => {
	it('renders the picker entry point', () => {
		const { body } = render(PlayOnDevice, {
			props: { mediaItemId: 'media-1', path: '/movie.mkv' }
		});
		expect(body).toContain('Play on device');
	});

	it('renders supported remote controls without enabling unsupported ones', () => {
		const { body } = renderWithTooltip(RendererRemote, {
			renderer: {
				id: 'uuid:tv',
				name: 'Living room',
				manufacturer: 'Test',
				modelName: 'TV',
				address: '192.168.1.2',
				lastSeen: new Date().toISOString(),
				expiresAt: new Date().toISOString(),
				supports: {
					play: true,
					pause: true,
					stop: true,
					seek: false,
					next: false,
					previous: false,
					volume: true,
					mute: true
				}
			}
		});
		expect(body).toContain('Living room');
		expect(body).toContain('aria-label="Play"');
		expect(body).toContain('Position');
		expect(body).toContain('disabled');
	});
});
