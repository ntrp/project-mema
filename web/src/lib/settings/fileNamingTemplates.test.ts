import { describe, expect, it } from 'vitest';
import { fileNamingTemplateExample, fileNamingTemplateSuggestions } from './fileNamingTemplates';

describe('file naming templates', () => {
	it('SCN-SETTINGS-007 ranks token suggestions and renders examples', () => {
		expect(
			fileNamingTemplateSuggestions('qual')
				.map((item) => item.param)
				.slice(0, 2)
		).toEqual(['quality', 'quality_full']);
		expect(fileNamingTemplateSuggestions('s00')[0]?.param).toBe('season:00');

		expect(
			fileNamingTemplateExample(
				'{movie_title} ({release_year}) S{season:00}E{episode:000} {quality_full}'
			)
		).toBe('The Matrix (1999) S01E003 Bluray-1080p Proper');
		expect(fileNamingTemplateExample('{unknown_token}')).toBe('{unknown_token}');
	});
});
