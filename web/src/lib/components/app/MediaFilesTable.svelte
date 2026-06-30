<script lang="ts">
	interface FileRow {
		kind: 'Media' | 'Metadata';
		path: string;
		fileName: string;
	}

	interface Props {
		filePaths: string[];
		metadataFilePaths: string[];
	}

	let { filePaths, metadataFilePaths }: Props = $props();

	const rows = $derived([
		...filePaths.map((path) => fileRow('Media', path)),
		...metadataFilePaths.map((path) => fileRow('Metadata', path))
	]);

	function fileRow(kind: FileRow['kind'], path: string): FileRow {
		return {
			kind,
			path,
			fileName: fileName(path)
		};
	}

	function fileName(path: string) {
		const normalized = path.replaceAll('\\', '/');
		return normalized.split('/').filter(Boolean).pop() ?? path;
	}
</script>

<section aria-labelledby="media-files-title">
	<h2 id="media-files-title">Files</h2>
	{#if rows.length}
		<div class="table-wrap media-files-table">
			<table>
				<thead>
					<tr>
						<th scope="col">File name</th>
						<th scope="col">Kind</th>
						<th scope="col">Path</th>
					</tr>
				</thead>
				<tbody>
					{#each rows as row (`${row.kind}:${row.path}`)}
						<tr>
							<td><strong>{row.fileName}</strong></td>
							<td>{row.kind}</td>
							<td><code>{row.path}</code></td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>
	{:else}
		<p class="empty">No imported files found.</p>
	{/if}
</section>
