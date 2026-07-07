import { describe, expect, it, vi } from 'vitest';

import MediaFilesHeader from '$lib/components/app/media/files/MediaFilesHeader.svelte';
import { renderWithTooltip } from '$lib/components/rendered/renderHelpers';
import type { MediaItem } from '$lib/settings/types';

describe('rendered rename controls', () => {
	it('shows a file overview rename action when files are available', () => {
		const { body } = renderWithTooltip(MediaFilesHeader, {
			item: {
				id: 'media-1',
				type: 'movie',
				title: 'Scenario Movie',
				filePaths: ['/library/Scenario Movie/Old.Name.mkv'],
				mediaFolderPath: '/library/Scenario Movie'
			} as MediaItem,
			canManage: true,
			scanningMediaItemId: undefined,
			onRename: vi.fn(),
			onRescanMediaFiles: vi.fn()
		});

		expect(body).toContain('Rename files');
		expect(body).toContain('Refresh file metadata');
		expect(body).toContain('file-pen-line');
	});
});
