const formatPatterns = [
	'Remux',
	'BluRay',
	'WEB-DL',
	'WEBDL',
	'WEBRip',
	'HDTV',
	'AMZN',
	'DSNP',
	'NF',
	'ATVP',
	'DD+',
	'TrueHD',
	'Atmos',
	'DTS',
	'HDR',
	'DV',
	'Proper',
	'Repack'
];

export function matchedFormats(value: string) {
	return formatPatterns.filter((format) => new RegExp(format.replace('+', '\\+'), 'i').test(value));
}
