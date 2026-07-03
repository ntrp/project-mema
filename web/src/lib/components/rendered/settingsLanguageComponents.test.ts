import { render } from 'svelte/server';
import { describe, expect, it, vi } from 'vitest';

import LanguageForm from '$lib/components/settings/languages/LanguageForm.svelte';
import LanguageTable from '$lib/components/settings/languages/LanguageTable.svelte';
import type { Language, LanguageForm as LanguageFormValue } from '$lib/settings/types';

describe('rendered language settings components (SCN-SETTINGS-003)', () => {
	it('renders language rows, aliases, deletion state, and empty state', () => {
		const populated = render(LanguageTable, {
			props: {
				languages: [language()],
				deletingCode: 'deu',
				onEdit: vi.fn(),
				onDelete: vi.fn()
			}
		});
		expect(populated.body).toContain('deu');
		expect(populated.body).toContain('German');
		expect(populated.body).toContain('DE, Deutsch');
		expect(populated.body).toContain('Deleting German');
		expect(populated.body).toContain('disabled');

		const empty = render(LanguageTable, {
			props: { languages: [], onEdit: vi.fn(), onDelete: vi.fn() }
		});
		expect(empty.body).toContain('No languages configured');
	});

	it('renders create and edit language form states', () => {
		const create = render(LanguageForm, {
			props: {
				form: languageForm(),
				saving: false,
				onSave: vi.fn(),
				onCancel: vi.fn()
			}
		});
		expect(create.body).toContain('ISO code');
		expect(create.body).toContain('Create language');
		expect(create.body).toContain('DE, DEU, GER, DEUTSCH');

		const edit = render(LanguageForm, {
			props: {
				form: languageForm({ originalCode: 'deu' }),
				saving: false,
				onSave: vi.fn(),
				onCancel: vi.fn()
			}
		});
		expect(edit.body).toContain('Update language');
		expect(edit.body).toContain('disabled');

		const saving = render(LanguageForm, {
			props: {
				form: languageForm({ originalCode: 'deu' }),
				saving: true,
				onSave: vi.fn(),
				onCancel: vi.fn()
			}
		});
		expect(saving.body).toContain('Saving');
	});
});

function language(overrides: Partial<Language> = {}): Language {
	return {
		code: 'deu',
		displayName: 'German',
		aliases: ['DE', 'Deutsch'],
		createdAt: '2026-07-03T00:00:00Z',
		updatedAt: '2026-07-03T00:00:00Z',
		...overrides
	};
}

function languageForm(overrides: Partial<LanguageFormValue> = {}): LanguageFormValue {
	return {
		code: 'deu',
		displayName: 'German',
		aliasesText: 'DE, Deutsch',
		...overrides
	};
}
