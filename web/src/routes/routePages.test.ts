import { describe, expect, it } from 'vitest';

import RootPage from './+page.svelte';
import ActivityPage from './activity/+page.svelte';
import ActivityBlocklistPage from './activity/blocklist/+page.svelte';
import ActivityHistoryPage from './activity/history/+page.svelte';
import BlacklistPage from './blacklist/+page.svelte';
import DiscoverPage from './discover/+page.svelte';
import MoviesPage from './movies/+page.svelte';
import RequestsPage from './requests/+page.svelte';
import SeriesPage from './series/+page.svelte';
import WantedPage from './wanted/+page.svelte';
import { renderPage } from './routeTestHelpers';

describe('top-level route pages (SCN-ROUTES-001)', () => {
	it('renders discover and library route sections from shell state', () => {
		const root = renderPage(RootPage);
		expect(root.body).toContain('Browse latest media');
		expect(root.body).toContain('No discovery sections available');

		const discover = renderPage(DiscoverPage);
		expect(discover.body).toContain('Browse latest media');

		const movies = renderPage(MoviesPage);
		expect(movies.body).toContain('Added movies');
		expect(movies.body).toContain('No movies added yet.');

		const series = renderPage(SeriesPage);
		expect(series.body).toContain('Added series');
		expect(series.body).toContain('No series added yet.');
	});

	it('renders activity, request, blacklist, and wanted fallbacks', () => {
		const activity = renderPage(ActivityPage);
		expect(activity.body).toContain('Activity queue');
		expect(activity.body).toContain('No queued activity');
		expect(activity.body).toContain('Queued downloads and in-progress grabs appear here');

		const activityHistory = renderPage(ActivityHistoryPage);
		expect(activityHistory.body).toContain('Activity history');
		expect(activityHistory.body).toContain('No completed background activity yet');
		expect(activityHistory.body).toContain('review finished media work');

		const activityBlocklist = renderPage(ActivityBlocklistPage);
		expect(activityBlocklist.body).toContain('Release blocklist');
		expect(activityBlocklist.body).toContain('No blocked releases yet');

		const requests = renderPage(RequestsPage);
		expect(requests.body).toContain('Media requests');
		expect(requests.body).toContain('No requests');

		const blacklist = renderPage(BlacklistPage);
		expect(blacklist.body).toContain('Blacklist');
		expect(blacklist.body).toContain('No blacklisted media');

		const wanted = renderPage(WantedPage);
		expect(wanted.body).toContain('Wanted');
		expect(wanted.body).toContain('No missing media.');
	});
});
