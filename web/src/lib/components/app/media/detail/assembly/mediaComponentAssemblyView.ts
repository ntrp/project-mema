import type {
	MediaComponentArtifact,
	MediaComponentAssemblyRun,
	MediaComponentSource,
	MediaItem
} from '$lib/settings/types';

export interface MediaComponentAssemblyView {
	retainedSources: MediaComponentSource[];
	baseSource?: MediaComponentSource;
	artifacts: MediaComponentArtifact[];
	allowedArtifacts: MediaComponentArtifact[];
	blockedCount: number;
	activeRun?: MediaComponentAssemblyRun;
	completedRun?: MediaComponentAssemblyRun;
	latestRun?: MediaComponentAssemblyRun;
	canAssemble: boolean;
	assembleLabel: string;
}

export function mediaComponentAssemblyView(item: MediaItem): MediaComponentAssemblyView {
	const retainedSources = (item.componentSources ?? []).filter(
		(source) => source.retentionState === 'retained'
	);
	const baseSource = retainedSources.find((source) => source.sourceRole === 'baseVideo');
	const artifacts = retainedSources.flatMap((source) => source.artifacts ?? []);
	const allowedArtifacts = baseSource
		? artifacts.filter(
				(artifact) => artifact.status === 'succeeded' && artifactAllowed(artifact, retainedSources)
			)
		: [];
	const assemblyRuns = [...(item.assemblyRuns ?? [])].sort((left, right) =>
		right.updatedAt.localeCompare(left.updatedAt)
	);
	const activeRun = assemblyRuns.find((run) => run.status === 'queued' || run.status === 'running');
	const completedRun = assemblyRuns.find((run) => run.status === 'succeeded');
	const latestRun = assemblyRuns[0];
	const canAssemble = Boolean(baseSource && allowedArtifacts.length > 0 && !activeRun);
	return {
		retainedSources,
		baseSource,
		artifacts,
		allowedArtifacts,
		blockedCount: blockedDecisionCount(retainedSources),
		activeRun,
		completedRun,
		latestRun,
		canAssemble,
		assembleLabel: latestRun?.status === 'failed' ? 'Retry assembly' : 'Start assembly'
	};
}

export function sourceDisplayName(source: MediaComponentSource) {
	return source.releaseTitle || fileName(source.sourceFilePath);
}

export function fileName(path: string) {
	return path.split('/').filter(Boolean).at(-1) ?? path;
}

export function sourceSummary(source: MediaComponentSource) {
	const artifactCount = source.artifacts?.length ?? 0;
	const decisions = source.compatibility ?? [];
	const blocked = decisions.some((decision) => decision.automationState === 'blocked');
	const suffix = artifactCount === 1 ? '1 artifact' : `${artifactCount} artifacts`;
	return blocked ? `${suffix}, review needed` : suffix;
}

export function statusTone(status: string) {
	if (status === 'succeeded' || status === 'retained' || status === 'allowed') return 'default';
	if (status === 'failed' || status === 'blocked' || status === 'rejected') return 'destructive';
	if (status === 'running' || status === 'queued' || status === 'pending') return 'secondary';
	return 'outline';
}

function artifactAllowed(artifact: MediaComponentArtifact, sources: MediaComponentSource[]) {
	const source = sources.find((item) => item.id === artifact.sourceId);
	return (source?.compatibility ?? []).some(
		(decision) =>
			decision.componentSourceId === artifact.sourceId &&
			(decision.automationState === 'allowed' || decision.reviewState === 'approved')
	);
}

function blockedDecisionCount(sources: MediaComponentSource[]) {
	return sources.reduce(
		(count, source) =>
			count +
			(source.compatibility ?? []).filter(
				(decision) =>
					decision.automationState === 'blocked' ||
					decision.reviewState === 'pending' ||
					decision.reviewState === 'rejected'
			).length,
		0
	);
}
