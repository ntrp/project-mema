import { describe, expect, it } from 'vitest';

import { parseArrCustomFormatImport } from './arrCustomFormatImport';

describe('arr custom format import (SCN-SETTINGS-009)', () => {
	it('imports wrapped formats with include and exclude specifications', () => {
		const [form] = parseArrCustomFormatImport(
			JSON.stringify({
				customFormats: [
					{
						name: 'Scenario Format',
						specifications: [
							{
								name: 'Required Source',
								implementation: 'SourceSpecification',
								required: true,
								fields: [{ name: 'source', value: 'WEB-DL' }]
							},
							{
								name: 'Rejected Language',
								implementationName: 'LanguageSpecification',
								negate: true,
								required: false,
								fields: [{ name: 'language', value: 'german' }]
							}
						]
					}
				]
			})
		);

		expect(form).toMatchObject({
			name: 'Scenario Format',
			includeInRenameTemplate: false,
			includeSpecs: [
				{
					id: 'required-source-web-dl',
					name: 'Required Source',
					type: 'source',
					value: 'WEB-DL',
					required: true
				}
			],
			excludeSpecs: [
				{
					id: 'rejected-language-german',
					name: 'Rejected Language',
					type: 'language',
					value: 'german',
					required: false
				}
			]
		});
	});

	it('supports single-object imports, direct values, and fallback field values', () => {
		const [form] = parseArrCustomFormatImport(
			JSON.stringify({
				name: 'Codec Rules',
				specifications: [
					{
						name: 'Video Codec',
						implementation: 'VideoCodecSpecification',
						value: 'x265'
					},
					{
						name: 'Any Term',
						implementation: 'ReleaseTitleSpecification',
						fields: [{ name: 'unused', value: true }]
					},
					{
						name: 'Ignored Empty',
						implementation: 'SourceSpecification',
						fields: [{ name: 'source', value: '' }]
					}
				]
			})
		);

		expect(form.includeSpecs.map((spec) => [spec.type, spec.value])).toEqual([
			['videoCodec', 'x265'],
			['releaseTitle', 'true']
		]);
	});

	it('reports invalid or unimportable payloads with user-facing errors', () => {
		expect(() => parseArrCustomFormatImport('[]')).toThrow('No custom formats found in JSON');
		expect(() => parseArrCustomFormatImport('{"name":"   ","specifications":[]}')).toThrow(
			'Imported custom format is missing a name'
		);
		expect(() =>
			parseArrCustomFormatImport('{"name":"Empty","specifications":[{"name":"No value"}]}')
		).toThrow('Empty does not contain importable specifications');
	});
});
