import js from '@eslint/js';
import tsParser from '@typescript-eslint/parser';
import tsPlugin from '@typescript-eslint/eslint-plugin';
import svelte from 'eslint-plugin-svelte';

const browserGlobals = {
	$derived: 'readonly',
	$state: 'readonly',
	console: 'readonly',
	document: 'readonly',
	Event: 'readonly',
	EventSource: 'readonly',
	HTMLDivElement: 'readonly',
	HTMLElement: 'readonly',
	HTMLSelectElement: 'readonly',
	KeyboardEvent: 'readonly',
	MessageEvent: 'readonly',
	SubmitEvent: 'readonly',
	URL: 'readonly',
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
		ignores: ['.svelte-kit/**', 'build/**', 'coverage/**', 'node_modules/**', 'test-results/**']
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
			]
		}
	}
];
