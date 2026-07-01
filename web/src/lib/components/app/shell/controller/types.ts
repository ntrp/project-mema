export interface ServerEventEnvelope<T = unknown> {
	id: string;
	type: string;
	time: string;
	data: T;
}

export type RelatedSectionKind = 'recommendations' | 'similar';
export type PeopleSectionKind = 'cast' | 'crew';
