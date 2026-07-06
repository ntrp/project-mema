import type { MediaComponentAssemblyRun, MediaComponentSource } from '$lib/settings/types';

export function componentSources(reviewState: 'pending' | 'approved'): MediaComponentSource[] {
	return [baseSource(), audioSource(reviewState)];
}

export function assemblyRun(
	overrides: Partial<MediaComponentAssemblyRun>
): MediaComponentAssemblyRun {
	return {
		id: 'run-1',
		mediaItemId: 'media-1',
		baseSourceId: 'source-base',
		outputPath: '/media/.mema/assemblies/run-1/assembled.mkv',
		status: 'queued',
		toolName: 'mkvmerge',
		toolSummary: '',
		createdAt: '2026-07-03T00:00:00Z',
		updatedAt: '2026-07-03T00:02:00Z',
		inputs: [
			{
				id: 'input-video',
				runId: 'run-1',
				sourceId: 'source-base',
				streamType: 'video',
				inputPath: '/media/.mema/components/base.mkv',
				provenance: { sourceFilePath: '/downloads/Base.Video.mkv' },
				createdAt: '2026-07-03T00:00:00Z'
			},
			{
				id: 'input-audio',
				runId: 'run-1',
				sourceId: 'source-audio',
				artifactId: 'artifact-audio',
				streamType: 'audio',
				inputPath: '/media/.mema/components/audio.mka',
				provenance: {
					sourceFilePath: '/downloads/Audio.Source.mkv',
					streamId: 1,
					language: 'jpn'
				},
				createdAt: '2026-07-03T00:00:00Z'
			}
		],
		...overrides
	};
}

function baseSource(): MediaComponentSource {
	return {
		id: 'source-base',
		mediaItemId: 'media-1',
		sourceRole: 'baseVideo',
		sourceFilePath: '/downloads/Base.Video.mkv',
		retainedPath: '/media/.mema/components/base.mkv',
		releaseTitle: 'Base.Video',
		streamInventory: '{}',
		retentionState: 'retained',
		retainedAt: '2026-07-03T00:00:00Z',
		createdAt: '2026-07-03T00:00:00Z',
		updatedAt: '2026-07-03T00:00:00Z'
	};
}

function audioSource(reviewState: 'pending' | 'approved'): MediaComponentSource {
	return {
		id: 'source-audio',
		mediaItemId: 'media-1',
		sourceRole: 'audio',
		sourceFilePath: '/downloads/Audio.Source.mkv',
		retainedPath: '/media/.mema/components/audio.mkv',
		releaseTitle: 'Audio.Source',
		streamInventory: '{}',
		retentionState: 'retained',
		retainedAt: '2026-07-03T00:00:00Z',
		createdAt: '2026-07-03T00:00:00Z',
		updatedAt: '2026-07-03T00:00:00Z',
		artifacts: [audioArtifact()],
		compatibility: [
			{
				id: 'decision-audio',
				mediaItemId: 'media-1',
				baseSourceId: 'source-base',
				componentSourceId: 'source-audio',
				confidenceState: reviewState === 'pending' ? 'uncertain' : 'likely',
				automationState: reviewState === 'pending' ? 'blocked' : 'allowed',
				reviewState,
				reason: 'runtime delta needs review',
				evidence: {},
				createdAt: '2026-07-03T00:00:00Z',
				updatedAt: '2026-07-03T00:00:00Z'
			}
		]
	};
}

function audioArtifact() {
	return {
		id: 'artifact-audio',
		mediaItemId: 'media-1',
		sourceId: 'source-audio',
		streamId: 1,
		streamType: 'audio' as const,
		language: 'jpn',
		outputPath: '/media/.mema/components/audio.mka',
		status: 'succeeded' as const,
		toolName: 'mkvextract',
		toolSummary: '',
		createdAt: '2026-07-03T00:00:00Z',
		updatedAt: '2026-07-03T00:00:00Z',
		completedAt: '2026-07-03T00:01:00Z'
	};
}
