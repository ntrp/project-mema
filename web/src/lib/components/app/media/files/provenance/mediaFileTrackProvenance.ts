import type { MediaFileRow } from '$lib/components/app/media/files/mediaFiles';

type MediaFileTrack = MediaFileRow['tracks'][number];
export type MediaFileTrackProvenance = NonNullable<MediaFileTrack['provenance']>;

interface ProvenanceField {
	label: string;
	value: string;
	multiline?: boolean;
}

export function provenanceFields(provenance: MediaFileTrackProvenance): ProvenanceField[] {
	return [
		field('ID', provenance.id),
		field('Media item', provenance.mediaItemId),
		field('Component type', provenance.componentType),
		field('Component key', provenance.componentKey),
		field('Release group', provenance.releaseGroup),
		field('Release name', provenance.releaseName),
		optionalField('Release ID', provenance.releaseId),
		optionalField('Source provider', provenance.sourceProvider),
		optionalField('Source file', provenance.sourceFilePath),
		optionalField('Retained source', provenance.retainedSourceId),
		optionalField('Source stream', provenance.sourceStreamId),
		field('Created', provenance.createdAt),
		field('Updated', provenance.updatedAt),
		field('Transformations', JSON.stringify(provenance.transformationChain, null, 2), true)
	].filter((item): item is ProvenanceField => Boolean(item));
}

function field(label: string, value: unknown, multiline = false): ProvenanceField {
	return { label, value: displayValue(value), multiline };
}

function optionalField(label: string, value: unknown): ProvenanceField | undefined {
	if (value === undefined || value === null || value === '') return undefined;
	return field(label, value);
}

function displayValue(value: unknown) {
	if (value === undefined || value === null || value === '') return '-';
	return String(value);
}
