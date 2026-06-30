const dateOnlyPattern = /^(\d{4})-(\d{2})-(\d{2})$/;

export function formatDate(value: string) {
	return formatBrowserDate(value, { dateStyle: 'medium' });
}

export function formatShortDate(value: string) {
	return formatBrowserDate(value, { dateStyle: 'short' });
}

export function formatDateTime(value: string) {
	return formatBrowserDate(value, { dateStyle: 'medium', timeStyle: 'short' });
}

export function formatShortDateTime(value: string) {
	return formatBrowserDate(value, { dateStyle: 'short', timeStyle: 'short' });
}

export function formatCompactDateTime(value: string) {
	return formatBrowserDate(value, {
		month: 'short',
		day: '2-digit',
		hour: '2-digit',
		minute: '2-digit'
	});
}

export function formatDateTimeWithSeconds(value: string) {
	return formatBrowserDate(value, {
		month: 'short',
		day: '2-digit',
		hour: '2-digit',
		minute: '2-digit',
		second: '2-digit'
	});
}

export function formatLongDateTime(value: string) {
	return formatBrowserDate(value, {
		year: 'numeric',
		month: 'short',
		day: '2-digit',
		hour: '2-digit',
		minute: '2-digit'
	});
}

export function formatTimeWithSeconds(value: string) {
	return formatBrowserDate(value, {
		hour: '2-digit',
		minute: '2-digit',
		second: '2-digit'
	});
}

function formatBrowserDate(value: string, options: Intl.DateTimeFormatOptions) {
	const date = parseDate(value);
	if (Number.isNaN(date.getTime())) {
		return value;
	}
	return new Intl.DateTimeFormat(browserLocales(), options).format(date);
}

function browserLocales() {
	const languages = globalThis.navigator?.languages;
	if (languages?.length) {
		return languages;
	}
	return globalThis.navigator?.language;
}

function parseDate(value: string) {
	const dateOnly = dateOnlyPattern.exec(value);
	if (!dateOnly) {
		return new Date(value);
	}
	return new Date(Number(dateOnly[1]), Number(dateOnly[2]) - 1, Number(dateOnly[3]));
}
