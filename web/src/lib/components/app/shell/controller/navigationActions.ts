import { goto } from '$app/navigation';
import { resolve } from '$app/paths';
import {
	settingsSectionHref,
	systemSectionHref
} from '$lib/components/app/navigation/appNavigation';
import type { HomeSection, SettingsSection, SystemSection } from '$lib/settings/types';
import type { AppShellState } from './state.svelte';

interface NavigationDeps {
	loadDiscoverSection: () => Promise<void>;
}

export function createNavigationActions(state: AppShellState, deps: NavigationDeps) {
	const loadDiscoverSection = deps.loadDiscoverSection;
	function selectHomeSection(section: HomeSection) {
		if (section === 'blacklist' && !state.isAdmin) {
			return;
		}
		state.activeView = 'home';
		state.activeHomeSection = section;
		void goto(resolve(`/${section}`));
	}

	function selectSettingsSection(section: string) {
		if (!state.isAdmin) {
			return;
		}
		state.activeSettingsSection = section as SettingsSection;
		void goto(resolve(settingsSectionHref(state.activeSettingsSection)));
	}

	function selectSystemSection(section: string) {
		if (!state.isAdmin) {
			return;
		}
		state.activeSystemSection = section as SystemSection;
		void goto(resolve(systemSectionHref(state.activeSystemSection)));
	}

	function selectSubmenuSection(section: string) {
		if (state.activeView === 'system') {
			selectSystemSection(section);
			return;
		}
		if (state.activePrimarySection === 'library') {
			selectHomeSection(section as HomeSection);
			return;
		}
		if (state.activePrimarySection === 'discover') {
			if (section === 'discover') {
				selectHomeSection('discover');
				return;
			}
			state.activeView = 'discover-section';
			state.activeHomeSection = 'discover';
			state.activeDiscoverSectionId = section;
			state.discoverSection = undefined;
			state.discoverSectionPage = 1;
			state.discoverSectionHasMore = true;
			void goto(resolve('/discover/[sectionId]', { sectionId: section }));
			void loadDiscoverSection();
			return;
		}
		selectSettingsSection(section);
	}

	function selectPrimarySection(section: string) {
		if (section === 'library') {
			selectHomeSection('movies');
			return;
		}
		if (section === 'settings') {
			if (!state.isAdmin) {
				return;
			}
			state.activeView = 'settings';
			state.activeSettingsSection = 'general';
			void goto(resolve('/settings/general'));
			return;
		}
		if (section === 'system') {
			if (!state.isAdmin) {
				return;
			}
			state.activeView = 'system';
			state.activeSystemSection = 'status';
			void goto(resolve('/system/status'));
			return;
		}
		selectHomeSection(section as HomeSection);
	}

	return {
		selectHomeSection,
		selectSettingsSection,
		selectSystemSection,
		selectSubmenuSection,
		selectPrimarySection
	};
}
