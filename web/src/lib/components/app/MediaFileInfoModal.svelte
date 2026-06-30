<script lang="ts">
	import type { MediaFileRow } from './mediaFiles';

	interface Props {
		row: MediaFileRow;
		onClose: () => void;
	}

	let { row, onClose }: Props = $props();

	const fields = $derived([
		['Relative path', row.relativePath],
		['Quality', row.quality],
		['Video codec', row.videoCodec],
		['Audio info', row.audioInfo],
		['Languages', row.languages],
		['Formats', row.formats.join(', ') || '-'],
		['Score', String(row.score)]
	]);
</script>

<div class="modal-backdrop" role="presentation" onclick={onClose}>
	<div
		class="modal-shell settings-modal media-file-modal"
		role="dialog"
		aria-modal="true"
		aria-labelledby="media-file-info-title"
		tabindex="-1"
		onclick={(event) => event.stopPropagation()}
		onkeydown={(event) => event.stopPropagation()}
	>
		<div class="modal-heading">
			<h2 id="media-file-info-title">File details</h2>
			<button type="button" class="icon-button" aria-label="Close" onclick={onClose}>
				<span class="app-icon" aria-hidden="true">close</span>
			</button>
		</div>
		<div class="metadata-facts">
			{#each fields as [label, value] (label)}
				<div>
					<span>{label}</span>
					<strong>{value}</strong>
				</div>
			{/each}
		</div>
		<div class="stream-detail-grid">
			<section>
				<h3>Video</h3>
				<p>Codec: {row.videoCodec}</p>
				<p>Resolution: {row.quality}</p>
				<p>Bitrate: -</p>
			</section>
			<section>
				<h3>Audio</h3>
				<p>Codec: {row.audioInfo}</p>
				<p>Languages: {row.languages}</p>
				<p>Channels: -</p>
			</section>
			<section>
				<h3>Subtitles</h3>
				<p>Languages: -</p>
				<p>Forced: -</p>
				<p>Count: -</p>
			</section>
		</div>
	</div>
</div>
