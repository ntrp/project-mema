<script lang="ts">
	import { client } from '$lib/api/client';

	type HealthState = 'idle' | 'loading' | 'ready' | 'error';

	let status: HealthState = $state('idle');
	let healthMessage = $state('Not checked yet');
	let errorMessage = $state('');

	async function checkHealth() {
		status = 'loading';
		errorMessage = '';

		try {
			const { data } = await client.GET('/health');

			if (!data) {
				throw new Error('API did not return health data');
			}

			status = 'ready';
			healthMessage = `${data.status} · ${data.version}`;
		} catch (error) {
			status = 'error';
			healthMessage = 'API unavailable';
			errorMessage = error instanceof Error ? error.message : 'Request failed';
		}
	}
</script>

<svelte:head>
	<title>Media Manager</title>
	<meta
		name="description"
		content="Self-hosted media manager for video libraries, profiles, indexers, downloads, and track assembly"
	/>
</svelte:head>

<main class="shell">
	<section class="overview" aria-labelledby="overview-title">
		<p class="eyebrow">Video-first media automation</p>
		<h1 id="overview-title">Media manager scaffold</h1>
		<p class="summary">
			Movies, TV, anime-specific metadata, subtitles, indexers, downloads, quality profiles, and
			track assembly are the first target surface.
		</p>

		<div class="actions">
			<button type="button" onclick={checkHealth} disabled={status === 'loading'}>
				{status === 'loading' ? 'Checking' : 'Check API'}
			</button>
		</div>
	</section>

	<section class="status-grid" aria-label="System status">
		<article>
			<span>API</span>
			<strong>{healthMessage}</strong>
			{#if errorMessage}
				<p>{errorMessage}</p>
			{/if}
		</article>

		<article>
			<span>Frontend</span>
			<strong>SvelteKit static</strong>
		</article>

		<article>
			<span>Contract</span>
			<strong>OpenAPI first</strong>
		</article>
	</section>
</main>

<style>
	:global(body) {
		margin: 0;
		font-family:
			Inter,
			ui-sans-serif,
			system-ui,
			-apple-system,
			BlinkMacSystemFont,
			'Segoe UI',
			sans-serif;
		background: #f6f7f9;
		color: #171a1f;
	}

	.shell {
		min-height: 100vh;
		display: grid;
		align-content: center;
		gap: 32px;
		padding: 48px;
		box-sizing: border-box;
	}

	.overview {
		max-width: 760px;
	}

	.eyebrow {
		margin: 0 0 12px;
		font-size: 13px;
		font-weight: 700;
		text-transform: uppercase;
		color: #556171;
	}

	h1 {
		margin: 0;
		font-size: 48px;
		line-height: 1.04;
		font-weight: 760;
	}

	.summary {
		max-width: 680px;
		margin: 20px 0 0;
		font-size: 18px;
		line-height: 1.55;
		color: #4b5563;
	}

	.actions {
		margin-top: 28px;
	}

	button {
		min-width: 132px;
		min-height: 40px;
		border: 1px solid #1f2937;
		border-radius: 6px;
		background: #1f2937;
		color: #ffffff;
		font: inherit;
		font-weight: 700;
		cursor: pointer;
	}

	button:disabled {
		cursor: wait;
		opacity: 0.7;
	}

	.status-grid {
		display: grid;
		grid-template-columns: repeat(3, minmax(0, 1fr));
		gap: 16px;
		max-width: 960px;
	}

	article {
		min-height: 108px;
		padding: 18px;
		border: 1px solid #d9dee7;
		border-radius: 8px;
		background: #ffffff;
		box-sizing: border-box;
	}

	article span {
		display: block;
		font-size: 13px;
		font-weight: 700;
		color: #6b7280;
	}

	article strong {
		display: block;
		margin-top: 10px;
		font-size: 19px;
	}

	article p {
		margin: 10px 0 0;
		color: #b42318;
	}

	@media (max-width: 760px) {
		.shell {
			padding: 28px;
		}

		h1 {
			font-size: 36px;
		}

		.status-grid {
			grid-template-columns: 1fr;
		}
	}
</style>
