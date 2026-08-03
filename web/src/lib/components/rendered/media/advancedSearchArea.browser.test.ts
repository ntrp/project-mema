import { mount, tick, unmount } from 'svelte';
import { afterEach, describe, expect, it, vi } from 'vitest';

import AdvancedSearchArea from '$lib/components/app/search/AdvancedSearchArea.svelte';

let component: ReturnType<typeof mount> | undefined;
let target = document.createElement('div');

afterEach(async () => {
	if (component) await unmount(component);
	target.remove();
	component = undefined;
});

describe('AdvancedSearchArea interactions', () => {
	it('waits for explicit submission before searching', async () => {
		const onSearch = vi.fn();
		target = document.createElement('div');
		document.body.append(target);
		component = mount(AdvancedSearchArea, {
			target,
			props: {
				initialQuery: 'scenario',
				metadataProviders: [],
				groups: [],
				searching: false,
				actionLabel: 'Add',
				onSearch,
				onAdd: vi.fn()
			}
		});
		await tick();

		expect(onSearch).not.toHaveBeenCalled();
		const button = target.querySelector('button[type="submit"]');
		if (!(button instanceof window.HTMLButtonElement)) throw new Error('missing search button');
		button.click();
		await tick();

		expect(onSearch).toHaveBeenCalledOnce();
		expect(onSearch).toHaveBeenCalledWith(expect.objectContaining({ query: 'scenario' }));
	});
});
