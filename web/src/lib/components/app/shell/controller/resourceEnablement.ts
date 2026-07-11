import type { AppShellState } from './state.svelte';

export function createResourceEnablement(state: AppShellState) {
	const settings = (section: string) =>
		state.authenticated &&
		state.isAdmin &&
		state.activeView === 'settings' &&
		state.activeSettingsSection === section;
	const home = (...sections: string[]) =>
		state.authenticated &&
		state.activeView === 'home' &&
		sections.includes(state.activeHomeSection);
	const mediaDetail = () => state.authenticated && Boolean(state.selectedMediaItemId);
	const mediaAction = () => state.authenticated && Boolean(state.activeMediaCandidate);
	const discovery = () =>
		state.authenticated &&
		state.isAdmin &&
		(home('discover', 'blacklist') ||
			['discover-section', 'discover-movies', 'discover-series', 'related-section'].includes(
				state.activeView
			));
	const system = (section: string) =>
		state.authenticated &&
		state.isAdmin &&
		state.activeView === 'system' &&
		state.activeSystemSection === section;

	return {
		languages: () =>
			settings('languages') || settings('profiles') || home('wanted') || mediaDetail(),
		tags: () => settings('tags') || settings('indexers') || home('requests') || mediaAction(),
		users: () => settings('users'),
		downloadClients: () => settings('download-clients'),
		indexers: () => settings('indexers'),
		metadataProviders: () =>
			settings('metadata') ||
			settings('library') ||
			(state.authenticated && state.isAdmin && state.activeView === 'advanced-search'),
		subtitleProviders: () => settings('subtitles'),
		libraryFolders: () => settings('library') || home('requests') || mediaDetail() || mediaAction(),
		pathMappings: () => settings('library'),
		mediaProfiles: () =>
			settings('profiles') ||
			settings('library') ||
			home('requests') ||
			mediaDetail() ||
			mediaAction(),
		customFormats: () => settings('custom-formats') || settings('profiles'),
		indexerSearch: () => settings('indexers') || system('indexing'),
		metadataCache: () => system('metadata'),
		profile: () => state.authenticated && state.activeView === 'profile',
		discovery
	};
}
