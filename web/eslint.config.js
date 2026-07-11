import js from '@eslint/js';
import tsParser from '@typescript-eslint/parser';
import tsPlugin from '@typescript-eslint/eslint-plugin';
import svelte from 'eslint-plugin-svelte';

const browserGlobals = {
	AbortController: 'readonly',
	AbortSignal: 'readonly',
	$derived: 'readonly',
	$state: 'readonly',
	console: 'readonly',
	document: 'readonly',
	Element: 'readonly',
	Event: 'readonly',
	EventSource: 'readonly',
	EventSourceInit: 'readonly',
	EventTarget: 'readonly',
	HTMLDivElement: 'readonly',
	HTMLElement: 'readonly',
	HTMLInputElement: 'readonly',
	HTMLSelectElement: 'readonly',
	HTMLTableRowElement: 'readonly',
	IntersectionObserver: 'readonly',
	KeyboardEvent: 'readonly',
	MessageEvent: 'readonly',
	MouseEvent: 'readonly',
	PointerEvent: 'readonly',
	SubmitEvent: 'readonly',
	URL: 'readonly',
	URLSearchParams: 'readonly',
	WheelEvent: 'readonly',
	window: 'readonly'
};

const nodeGlobals = {
	process: 'readonly',
	URL: 'readonly'
};

const domTypeGlobals = {
	FileList: 'readonly',
	HTMLElement: 'readonly',
	HTMLAnchorElement: 'readonly',
	HTMLButtonElement: 'readonly',
	HTMLInputElement: 'readonly',
	HTMLParagraphElement: 'readonly',
	HTMLSpanElement: 'readonly',
	HTMLTableRowElement: 'readonly',
	HTMLTableSectionElement: 'readonly'
};

export default [
	js.configs.recommended,
	...svelte.configs['flat/recommended'],
	{
		ignores: [
			'.svelte-kit/**',
			'build/**',
			'coverage/**',
			'node_modules/**',
			'src/lib/api/generated/**',
			'test-results/**'
		]
	},
	{
		files: ['**/*.{js,ts}'],
		languageOptions: {
			parser: tsParser,
			parserOptions: {
				sourceType: 'module'
			},
			globals: browserGlobals
		},
		plugins: {
			'@typescript-eslint': tsPlugin
		},
		rules: {
			...tsPlugin.configs.recommended.rules,
			'no-restricted-syntax': [
				'error',
				{
					selector: "NewExpression[callee.name='EventSource']",
					message:
						'Use the application SSE transport in src/lib/app/realtime/appEventSource.ts instead of opening another EventSource.'
				}
			],
			'no-unused-vars': 'off',
			'@typescript-eslint/no-unused-vars': [
				'error',
				{ argsIgnorePattern: '^_', varsIgnorePattern: '^_' }
			]
		}
	},
	{
		// These streams are scoped interactive/diagnostic transports, not the
		// application-wide notification stream served by /api/events.
		files: [
			'src/lib/app/realtime/appEventSource.ts',
			'src/lib/components/app/media/release-search/releaseSearchStream.ts',
			'src/lib/components/settings/system/logs/SystemLogsSettings.svelte'
		],
		rules: {
			'no-restricted-syntax': 'off'
		}
	},
	{
		files: ['playwright.config.ts', 'e2e/**/*.{js,mjs,ts}'],
		languageOptions: {
			globals: nodeGlobals
		}
	},
	{
		files: ['src/lib/components/ui/**/*.{svelte,ts}', 'src/lib/utils.ts'],
		languageOptions: {
			globals: domTypeGlobals
		}
	},
	{
		files: ['**/*.svelte'],
		languageOptions: {
			parserOptions: {
				parser: tsParser
			},
			globals: browserGlobals
		},
		plugins: {
			'@typescript-eslint': tsPlugin
		},
		rules: {
			'no-restricted-syntax': [
				'error',
				{
					selector: "NewExpression[callee.name='EventSource']",
					message:
						'Use the application SSE transport in src/lib/app/realtime/appEventSource.ts instead of opening another EventSource.'
				}
			],
			'no-unused-vars': 'off',
			'@typescript-eslint/no-unused-vars': [
				'error',
				{ argsIgnorePattern: '^_', varsIgnorePattern: '^_' }
			],
			'svelte/no-navigation-without-resolve': 'off'
		}
	},
	{
		files: ['src/lib/components/settings/system/logs/SystemLogsSettings.svelte'],
		rules: {
			'no-restricted-syntax': 'off'
		}
	}
];
