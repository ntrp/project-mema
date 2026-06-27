<script lang="ts">
	import SidebarMenu from './SidebarMenu.svelte';
	import type { HomeSection } from '$lib/settings/types';

	interface Props {
		activeSection: HomeSection;
		onSelect: (_section: HomeSection) => void;
	}

	let { activeSection, onSelect }: Props = $props();

	const homeItems = [
		{ value: 'explore', label: 'Explore', meta: 'Digest' },
		{ value: 'movies', label: 'Movies', meta: 'Anime included' },
		{ value: 'series', label: 'Series', meta: 'Anime included' },
		{ value: 'activity', label: 'Activity', meta: 'Queue' }
	] satisfies { value: HomeSection; label: string; meta: string }[];

	const latest = [
		{ title: 'Dune: Part Two', type: 'Movie', year: '2024', status: 'Trending' },
		{ title: "Frieren: Beyond Journey's End", type: 'Series', year: '2023', status: 'Anime' },
		{ title: 'The Last of Us', type: 'Series', year: '2023', status: 'Popular' },
		{ title: 'Suzume', type: 'Movie', year: '2022', status: 'Anime' }
	];
	const movies = [
		{ title: 'Blade Runner 2049', quality: '2160p HDR', state: 'Monitored' },
		{ title: 'Your Name', quality: '1080p', state: 'Monitored' },
		{ title: 'Dune', quality: '2160p', state: 'Available' }
	];
	const series = [
		{ title: 'Attack on Titan', seasons: '4 seasons', state: 'Complete' },
		{ title: 'Severance', seasons: '2 seasons', state: 'Monitoring' },
		{ title: 'Arcane', seasons: '2 seasons', state: 'Monitoring' }
	];
	const activity = [
		{ title: 'The Apothecary Diaries S02E03', progress: '72%', state: 'Downloading' },
		{ title: 'Dune: Part Two', progress: 'Waiting', state: 'Import queue' },
		{ title: 'Frieren S01E12', progress: 'Syncing', state: 'Post-processing' }
	];
</script>

<div class="workspace-layout">
	<SidebarMenu
		title="Library"
		items={homeItems}
		active={activeSection}
		onSelect={(section) => onSelect(section as HomeSection)}
	/>

	<section class="workspace-main" aria-labelledby="home-title">
		{#if activeSection === 'explore'}
			<div class="page-heading">
				<p>Explore</p>
				<h1 id="home-title">Latest media digest</h1>
			</div>
			<div class="digest-grid">
				{#each latest as item (item.title)}
					<article class="media-tile">
						<div class="poster-placeholder">{item.type}</div>
						<h2>{item.title}</h2>
						<p>{item.year} · {item.status}</p>
					</article>
				{/each}
			</div>
		{:else if activeSection === 'movies'}
			<div class="page-heading">
				<p>Movies</p>
				<h1 id="home-title">Added movies</h1>
			</div>
			<div class="data-list">
				{#each movies as item (item.title)}
					<div class="data-row">
						<strong>{item.title}</strong>
						<span>{item.quality}</span>
						<small>{item.state}</small>
					</div>
				{/each}
			</div>
		{:else if activeSection === 'series'}
			<div class="page-heading">
				<p>Series</p>
				<h1 id="home-title">Added series</h1>
			</div>
			<div class="data-list">
				{#each series as item (item.title)}
					<div class="data-row">
						<strong>{item.title}</strong>
						<span>{item.seasons}</span>
						<small>{item.state}</small>
					</div>
				{/each}
			</div>
		{:else}
			<div class="page-heading">
				<p>Activity</p>
				<h1 id="home-title">Downloads and imports</h1>
			</div>
			<div class="data-list">
				{#each activity as item (item.title)}
					<div class="data-row">
						<strong>{item.title}</strong>
						<span>{item.progress}</span>
						<small>{item.state}</small>
					</div>
				{/each}
			</div>
		{/if}
	</section>
</div>
