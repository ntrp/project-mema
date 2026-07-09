import js from '@eslint/js';
import tsParser from '@typescript-eslint/parser';
import tsPlugin from '@typescript-eslint/eslint-plugin';
import svelte from 'eslint-plugin-svelte';

const browserGlobals = {
	$derived: 'readonly',
	$state: 'readonly',
	console: 'readonly',
	document: 'readonly',
	Element: 'readonly',
	Event: 'readonly',
	EventSource: 'readonly',
	EventTarget: 'readonly',
	HTMLDivElement: 'readonly',
	HTMLElement: 'readonly',
	HTMLInputElement: 'readonly',
	HTMLSelectElement: 'readonly',
	HTMLTableRowElement: 'readonly',
	IntersectionObserver: 'readonly',
	KeyboardEvent: 'readonly',
	MessageEvent: 'readonly',
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
			'no-unused-vars': 'off',
			'@typescript-eslint/no-unused-vars': [
				'error',
				{ argsIgnorePattern: '^_', varsIgnorePattern: '^_' }
			]
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
			'no-unused-vars': 'off',
			'@typescript-eslint/no-unused-vars': [
				'error',
				{ argsIgnorePattern: '^_', varsIgnorePattern: '^_' }
			],
			'svelte/no-navigation-without-resolve': 'off'
		}
	}
];
