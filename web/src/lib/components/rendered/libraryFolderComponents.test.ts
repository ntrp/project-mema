import { render } from 'svelte/server';
import { beforeEach, describe, expect, it, vi } from 'vitest';

import TemplateTokenTextarea from '$lib/components/settings/library/TemplateTokenTextarea.svelte';
import LibraryFolderCreateForm from '$lib/components/settings/library/folders/LibraryFolderCreateForm.svelte';
import LibraryFolderTree from '$lib/components/settings/library/folders/LibraryFolderTree.svelte';
import {
	createNode,
	flattenTree,
	type TreeNode
} from '$lib/components/settings/library/folders/libraryFolderTree';

const apiMock = vi.hoisted(() => ({
	listLibraryFolderOptions: vi.fn(),
	createLibraryFolderOption: vi.fn()
}));

vi.mock('$lib/settings/api', async (importOriginal) => ({
	...(await importOriginal<typeof import('$lib/settings/api')>()),
	listLibraryFolderOptions: apiMock.listLibraryFolderOptions,
	createLibraryFolderOption: apiMock.createLibraryFolderOption
}));

describe('rendered library folder controls (SCN-SETTINGS-016)', () => {
	beforeEach(() => {
		apiMock.listLibraryFolderOptions.mockReset();
		apiMock.createLibraryFolderOption.mockReset();
	});

	it('renders create-folder state and disables invalid submissions', () => {
		const empty = render(LibraryFolderCreateForm, {
			props: { name: '   ', disabled: false, creating: false, onCreate: vi.fn() }
		});
		expect(empty.body).toContain('Create under selected folder');
		expect(empty.body).toContain('New folder name');
		expect(empty.body).toContain('disabled');

		const creating = render(LibraryFolderCreateForm, {
			props: { name: 'Movies', disabled: false, creating: true, onCreate: vi.fn() }
		});
		expect(creating.body).toContain('Creating');
	});

	it('renders folder tree selection, expansion, loading, and errors', () => {
		const { body } = render(LibraryFolderTree, {
			props: {
				visibleNodes: [
					treeNode({
						name: 'media',
						path: '/media',
						expanded: true,
						loaded: true,
						children: ['/media/movies']
					}),
					treeNode({
						name: 'movies',
						path: '/media/movies',
						parentPath: '/media',
						depth: 1,
						loading: true,
						error: 'Permission denied'
					})
				],
				selectedPath: '/media/movies',
				onToggle: vi.fn()
			}
		});

		expect(body).toContain('Server folders');
		expect(body).toContain('/media/movies');
		expect(body).toContain('Permission denied');
		expect(body).toContain('aria-selected="true"');
		expect(body).toContain('...');
	});

	it('flattens only expanded folder descendants', () => {
		const root = treeNode({ path: '/media', children: ['/media/movies'], expanded: true });
		const movies = treeNode({
			name: 'movies',
			path: '/media/movies',
			parentPath: '/media',
			depth: 1,
			children: ['/media/movies/action'],
			expanded: false
		});
		const action = treeNode({
			name: 'action',
			path: '/media/movies/action',
			parentPath: '/media/movies',
			depth: 2
		});

		expect(
			flattenTree({ [root.path]: root, [movies.path]: movies, [action.path]: action }, root.path)
		).toEqual([root, movies]);
		expect(createNode({ name: 'shows', path: '/media/shows' }, '/media', 1)).toMatchObject({
			name: 'shows',
			parentPath: '/media',
			expanded: false,
			loaded: false
		});
	});
});

describe('rendered file naming token input (SCN-SETTINGS-007)', () => {
	it('renders the current template value as editable text', () => {
		const { body } = render(TemplateTokenTextarea, {
			props: {
				value: '{Media Title} ({Release Year})/{Quality Title}',
				onChange: vi.fn()
			}
		});

		expect(body).toContain('{Media Title} ({Release Year})/{Quality Title}');
		expect(body).toContain('font-mono');
	});
});

function treeNode(overrides: Partial<TreeNode>): TreeNode {
	return {
		name: 'media',
		path: '/media',
		depth: 0,
		expanded: false,
		loaded: false,
		loading: false,
		error: '',
		children: [],
		...overrides
	};
}
