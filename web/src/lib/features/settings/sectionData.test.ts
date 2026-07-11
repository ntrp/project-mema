import { describe, expect, it } from 'vitest';
import { loadMediaActionSettings, loadSettingsSection, loadSystemSettings } from './sectionData';

describe('legacy settings loaders', () => {
	it('defers settings, system, and media resources to feature queries', async () => {
		await expect(loadSettingsSection('indexers')).resolves.toEqual({});
		await expect(loadSystemSettings('metadata')).resolves.toEqual({});
		await expect(loadMediaActionSettings()).resolves.toEqual({});
	});
});
