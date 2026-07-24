import type { DLNARendererProfile } from '$lib/settings/types';

export interface DLNAProfileTableProps {
	profiles: DLNARendererProfile[];
	search: string;
	selectedId?: string;
	onSearch: (_value: string) => void;
	onEdit: (_profile: DLNARendererProfile) => void;
	onClone: (_profile: DLNARendererProfile) => void;
	onReset: (_profile: DLNARendererProfile) => void;
	onExport: (_profile: DLNARendererProfile) => void;
	onDelete: (_profile: DLNARendererProfile) => void;
	onCreate?: () => void;
	onImport?: () => void;
	onTrace?: () => void;
	onRestoreOriginals?: () => void;
}
