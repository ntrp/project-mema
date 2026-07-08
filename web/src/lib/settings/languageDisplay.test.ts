import { describe, expect, it } from 'vitest';
import { displayLanguage, languageMatchKey } from './languageDisplay';

describe('language display', () => {
	it('SCN-SETTINGS-003 normalizes common language aliases for display and matching', () => {
		expect(languageMatchKey(' German Language ')).toBe('de');
		expect(languageMatchKey('ita')).toBe('it');
		expect(languageMatchKey('jpn')).toBe('ja');
		expect(languageMatchKey('Portuguese')).toBe('pt');
		expect(displayLanguage('-')).toBe('-');
		expect(displayLanguage('english')).toMatch(/English/i);
	});
});
