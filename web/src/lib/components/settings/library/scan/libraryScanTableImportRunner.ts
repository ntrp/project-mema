import { importCheckedScanRows } from './libraryScanTableActions';
import type { MatchDraft } from './libraryScanImport';
import type { ImportBulkOptions } from './libraryScanImportPayloads';
import type { LibraryScan, LibraryScanImportRequest, LibraryScanItem } from '$lib/settings/types';

export async function runCheckedScanImport(input: {
	canImport: boolean;
	checkedRows: LibraryScanItem[];
	allRows: LibraryScanItem[];
	drafts: Record<string, MatchDraft>;
	bulk: ImportBulkOptions;
	scan: LibraryScan;
	setImporting: (_value: boolean) => void;
	setImportingItem: (_id: string) => void;
	onImport: (_scan: LibraryScan, _request: LibraryScanImportRequest) => Promise<void>;
}) {
	if (!input.canImport) return;
	input.setImporting(true);
	try {
		await importCheckedScanRows({
			canImport: input.canImport,
			checkedRows: input.checkedRows,
			allRows: input.allRows,
			drafts: input.drafts,
			scan: input.scan,
			onProgress: input.setImportingItem,
			onImport: input.onImport,
			bulk: input.bulk
		});
	} finally {
		input.setImportingItem('');
		input.setImporting(false);
	}
}
