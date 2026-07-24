import { page } from 'vitest/browser';
import { mount, unmount } from 'svelte';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';

const api = vi.hoisted(() => ({
	listRenderers: vi.fn(),
	playOnRenderer: vi.fn(),
	rendererStatus: vi.fn(),
	rendererTransport: vi.fn(),
	seekRenderer: vi.fn(),
	setRendererMute: vi.fn(),
	setRendererVolume: vi.fn()
}));
vi.mock('./api', () => api);

import RenderWithTooltip from '$lib/components/rendered/RenderWithTooltip.svelte';
import PlayOnDevice from './PlayOnDevice.svelte';
import RendererRemote from './RendererRemote.svelte';
import type { Renderer } from './types';

const renderer: Renderer = {
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
		seek: true,
		next: false,
		previous: false,
		volume: true,
		mute: true
	}
};

const browserDescribe = typeof document === 'undefined' ? describe.skip : describe;

browserDescribe('renderer component interactions', () => {
	let components: object[] = [];

	beforeEach(() => {
		vi.clearAllMocks();
		document.body.replaceChildren();
		api.rendererStatus.mockResolvedValue({
			transportState: 'PLAYING',
			transportStatus: 'OK',
			positionSeconds: 12,
			durationSeconds: 120,
			volume: 40,
			muted: false
		});
		api.rendererTransport.mockResolvedValue(undefined);
		api.seekRenderer.mockResolvedValue(undefined);
		api.setRendererMute.mockResolvedValue(undefined);
		api.setRendererVolume.mockResolvedValue(undefined);
	});

	afterEach(async () => {
		for (const component of components) await unmount(component);
		components = [];
		vi.restoreAllMocks();
	});

	it('polls, cleans up its timer, and runs every remote interaction', async () => {
		let poll = () => {};
		const clearInterval = vi.spyOn(globalThis, 'clearInterval').mockImplementation(() => {});
		vi.spyOn(globalThis, 'setInterval').mockImplementation(
			(handler: Parameters<typeof globalThis.setInterval>[0]) => {
				poll = handler as () => void;
				return 17 as unknown as ReturnType<typeof globalThis.setInterval>;
			}
		);
		const component = mountWithTooltip(RendererRemote, { renderer });
		components.push(component);
		await vi.waitFor(() => expect(api.rendererStatus).toHaveBeenCalledTimes(1));
		poll();
		await vi.waitFor(() => expect(api.rendererStatus).toHaveBeenCalledTimes(2));

		button('Play').click();
		await vi.waitFor(() => expect(api.rendererTransport).toHaveBeenCalledWith('uuid:tv', 'play'));
		const ranges = [...document.querySelectorAll<HTMLInputElement>('input[type="range"]')];
		ranges[0].value = '45';
		ranges[0].dispatchEvent(new Event('change', { bubbles: true }));
		await vi.waitFor(() => expect(api.seekRenderer).toHaveBeenCalledWith('uuid:tv', 45));
		await vi.waitFor(() => expect(ranges[1].disabled).toBe(false));
		ranges[1].value = '55';
		ranges[1].dispatchEvent(new Event('change', { bubbles: true }));
		await vi.waitFor(() => expect(api.setRendererVolume).toHaveBeenCalledWith('uuid:tv', 55));
		await vi.waitFor(() => expect(button('Mute').disabled).toBe(false));
		button('Mute').click();
		await vi.waitFor(() => expect(api.setRendererMute).toHaveBeenCalledWith('uuid:tv', true));

		await unmount(component);
		components = [];
		expect(clearInterval).toHaveBeenCalledWith(17);
	});

	it('renders command errors and keeps unsupported controls disabled', async () => {
		api.rendererTransport.mockRejectedValueOnce(new Error('Transition unavailable'));
		const unsupported = {
			...renderer,
			supports: { ...renderer.supports, pause: false, seek: false, volume: false, mute: false }
		};
		const component = mountWithTooltip(RendererRemote, { renderer: unsupported });
		components.push(component);
		await vi.waitFor(() => expect(api.rendererStatus).toHaveBeenCalled());
		expect(button('Pause').disabled).toBe(true);
		expect(button('Mute').disabled).toBe(true);
		expect(
			[...document.querySelectorAll<HTMLInputElement>('input[type="range"]')].every(
				(input) => input.disabled
			)
		).toBe(true);
		await page.getByText(/Position/).hover();
		await vi.waitFor(() => expect(document.body.textContent).toContain('Seek is not supported'));
		await page.getByText(/Volume/).hover();
		await vi.waitFor(() => expect(document.body.textContent).toContain('Volume is not supported'));
		button('Play').click();
		await vi.waitFor(() =>
			expect(document.querySelector('[role="alert"]')?.textContent).toContain(
				'Transition unavailable'
			)
		);
	});

	it('discovers a renderer, starts playback, and reports discovery errors', async () => {
		api.listRenderers.mockResolvedValueOnce([renderer]);
		api.playOnRenderer.mockResolvedValue(undefined);
		let component = mountWithTooltip(PlayOnDevice, {
			mediaItemId: '00000000-0000-0000-0000-000000000001',
			path: '/movie.mkv'
		});
		components.push(component);
		await page.getByRole('button', { name: 'Play on device' }).click();
		await vi.waitFor(() => expect(api.listRenderers).toHaveBeenCalledTimes(1));
		await vi.waitFor(() => expect(document.body.textContent).toContain('Living room'));
		await page.getByRole('button', { name: /Living room/ }).click();
		await vi.waitFor(() =>
			expect(api.playOnRenderer).toHaveBeenCalledWith(
				'uuid:tv',
				'00000000-0000-0000-0000-000000000001',
				'/movie.mkv'
			)
		);

		await unmount(component);
		components = [];
		document.body.replaceChildren();
		api.listRenderers.mockRejectedValueOnce(new Error('Discovery unavailable'));
		component = mountWithTooltip(PlayOnDevice, {
			mediaItemId: '00000000-0000-0000-0000-000000000001',
			path: '/movie.mkv'
		});
		components.push(component);
		await page.getByRole('button', { name: 'Play on device' }).click();
		await vi.waitFor(() =>
			expect(document.querySelector('[role="alert"]')?.textContent).toContain(
				'Discovery unavailable'
			)
		);
	});
});

function mountWithTooltip(component: unknown, componentProps: Record<string, unknown>) {
	return mount(RenderWithTooltip, {
		target: document.body,
		props: { component: component as never, componentProps }
	});
}

function button(label: string) {
	const found = [...document.querySelectorAll('button')].find(
		(item) => item.getAttribute('aria-label') === label || item.textContent?.trim() === label
	);
	if (!found) throw new Error(`Button ${label} not found`);
	return found;
}
