import { mount, unmount } from 'svelte';
import { afterEach, describe, expect, it, vi } from 'vitest';
import LibraryFolderPickerTestHost from './LibraryFolderPickerTestHost.svelte';
import LibraryFolderTree from './LibraryFolderTree.svelte';

const apiMock = vi.hoisted(() => ({
	listLibraryFolderOptions: vi.fn(),
	createLibraryFolderOption: vi.fn()
}));

vi.mock('$lib/settings/api', () => ({
	listLibraryFolderOptions: apiMock.listLibraryFolderOptions,
	createLibraryFolderOption: apiMock.createLibraryFolderOption
}));

const browserDescribe = typeof document === 'undefined' ? describe.skip : describe;

browserDescribe('LibraryFolderTree', () => {
	let component: object | undefined;

	afterEach(async () => {
		if (component) await unmount(component);
		component = undefined;
		document.body.replaceChildren();
		vi.clearAllMocks();
	});

	it('selects a leaf folder when its row is clicked', () => {
		const onSelect = vi.fn();
		component = mount(LibraryFolderTree, {
			target: document.body,
			props: {
				visibleNodes: [
					{
						name: 'movies',
						path: '/media/movies',
						depth: 0,
						expanded: false,
						loaded: true,
						loading: false,
						error: '',
						children: []
					}
				],
				selectedPath: '/',
				onSelect,
				onToggle: vi.fn()
			}
		});

		folderRow('/media/movies').click();
		expect(onSelect).toHaveBeenCalledWith('/media/movies');
	});

	it('shows root, preopens the main data path, lazily expands, selects, and creates folders', async () => {
		let created = false;
		apiMock.listLibraryFolderOptions.mockImplementation(async (path?: string) => {
			if (!path) return folderList('/data/media', [{ name: 'movies', path: '/data/media/movies' }]);
			if (path === '/') return folderList('/', [{ name: 'data', path: '/data' }]);
			if (path === '/data') return folderList('/data', [{ name: 'media', path: '/data/media' }]);
			if (path === '/data/media/movies') {
				return folderList(path, created ? [{ name: 'New', path: `${path}/New` }] : []);
			}
			throw new Error(`Unexpected path ${path}`);
		});
		apiMock.createLibraryFolderOption.mockImplementation(async (...args: [string, string]) => {
			const [parentPath, folderName] = args;
			created = true;
			return { name: folderName, path: `${parentPath}/${folderName}` };
		});

		component = mount(LibraryFolderPickerTestHost, {
			target: document.body,
			props: { onClose: vi.fn(), onUse: vi.fn() }
		});

		await vi.waitFor(() => expect(apiMock.listLibraryFolderOptions).toHaveBeenCalledTimes(3));
		expect(apiMock.listLibraryFolderOptions).toHaveBeenNthCalledWith(1, undefined);
		expect(apiMock.listLibraryFolderOptions).toHaveBeenNthCalledWith(2, '/');
		expect(apiMock.listLibraryFolderOptions).toHaveBeenNthCalledWith(3, '/data');
		expect(document.querySelector('[role="treeitem"]')?.textContent?.trim()).toBe('/');
		expect(selectedPath()).toBe('/data/media');

		folderExpandButton('/data/media/movies').click();
		await vi.waitFor(() =>
			expect(apiMock.listLibraryFolderOptions).toHaveBeenCalledWith('/data/media/movies')
		);
		folderRow('/data/media/movies').click();
		await vi.waitFor(() => expect(selectedPath()).toBe('/data/media/movies'));

		const input = document.querySelector('input[placeholder="New folder name"]');
		if (!input || !('value' in input)) throw new Error('New folder input not found');
		input.value = 'New';
		input.dispatchEvent(new Event('input', { bubbles: true }));
		await vi.waitFor(() => expect(createButton().disabled).toBe(false));
		expect(createButton().querySelector('svg')).not.toBeNull();
		createButton().click();

		await vi.waitFor(() =>
			expect(apiMock.createLibraryFolderOption).toHaveBeenCalledWith('/data/media/movies', 'New')
		);
		await vi.waitFor(() => expect(selectedPath()).toBe('/data/media/movies/New'));
	});
});

function folderList(currentPath: string, entries: { name: string; path: string }[]) {
	return { currentPath, entries };
}

function folderItem(path: string) {
	const item = [...document.querySelectorAll('[role="treeitem"]')].find((element) =>
		element.textContent?.includes(path)
	);
	if (!item) throw new Error(`Folder ${path} not found`);
	return item;
}

function folderExpandButton(path: string) {
	const button = folderItem(path).querySelector('button');
	if (!button) throw new Error(`Expand button for ${path} not found`);
	return button;
}

function folderRow(path: string) {
	const buttons = folderItem(path).querySelectorAll('button');
	const button = buttons[1] ?? buttons[0];
	if (!button) throw new Error(`Folder row for ${path} not found`);
	return button;
}

function createButton() {
	const button = [...document.querySelectorAll('button')].find((element) =>
		element.textContent?.includes('Create folder')
	);
	if (!button) throw new Error('Create folder button not found');
	return button;
}

function selectedPath() {
	return [...document.querySelectorAll('strong')].find((element) =>
		element.textContent?.includes('/data')
	)?.textContent;
}
