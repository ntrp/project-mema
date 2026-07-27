import type { MatchDraft } from './libraryScanImport';
import {
	duplicateRemovalPathsForRows,
	importPayloadForSingleRow,
	type ImportBulkOptions
} from './libraryScanImportPayloads';
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
		const duplicatePaths = duplicateRemovalPathsForRows(input.allRows, input.drafts);
		for (const [index, row] of input.checkedRows.entries()) {
			input.setImportingItem(row.id);
			await input.onImport(
				input.scan,
				importPayloadForSingleRow(
					row,
					input.drafts,
					input.bulk,
					index === 0 ? duplicatePaths : [],
					input.allRows
				)
			);
		}
	} finally {
		input.setImportingItem('');
		input.setImporting(false);
	}
}
