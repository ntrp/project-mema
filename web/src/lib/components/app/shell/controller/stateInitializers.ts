import type { AppRouteState } from './routeState';

export function initialiseRouteState(target: object, route: AppRouteState) {
	Object.assign(target, {
		route,
		activeView: route.view,
		activeHomeSection: route.homeSection,
		activeActivitySection: route.activitySection,
		activeSettingsSection: route.settingsSection,
		activeSystemSection: route.systemSection,
		activeDiscoverSectionId: route.discoverSectionId,
		activeDiscoverSubmenuSection: route.discoverSubmenuSection,
		activeRelatedSectionKind: route.relatedSectionKind,
		activePeopleSectionKind: route.peopleSectionKind,
		selectedMediaItemId: route.selectedMediaItemId,
		selectedRequestId: route.selectedRequestId,
		searchQuery: route.advancedQuery
	});
}
