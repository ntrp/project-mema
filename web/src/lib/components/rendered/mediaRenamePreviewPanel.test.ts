import { render } from 'svelte/server';
import { describe, expect, it, vi } from 'vitest';

import MediaRenamePreviewPanel from '$lib/components/app/media/files/MediaRenamePreviewPanel.svelte';

describe('rendered rename preview panel', () => {
	it('shows safe, unchanged, skipped, and blocked preview states', () => {
		const { body } = render(MediaRenamePreviewPanel, {
			props: {
				rows: [
					{
						currentPath: '/library/old.mkv',
						proposedPath: '/library/new.mkv',
						status: 'safe',
						messages: []
					},
					{
						currentPath: '/library/current.mkv',
						proposedPath: '/library/current.mkv',
						status: 'unchanged',
						messages: []
					},
					{
						currentPath: '/library/missing.mkv',
						proposedPath: '',
						status: 'missing',
						messages: ['File is missing.']
					},
					{
						currentPath: '/library/blocked.mkv',
						proposedPath: '',
						status: 'blocked',
						messages: ['Season and episode could not be detected.']
					}
				],
				loading: false,
				onPreview: vi.fn()
			}
		});

		expect(body).toContain('Rename preview');
		expect(body).toContain('safe');
		expect(body).toContain('unchanged');
		expect(body).toContain('skipped');
		expect(body).toContain('blocked');
		expect(body).toContain('Season and episode could not be detected.');
	});
});
