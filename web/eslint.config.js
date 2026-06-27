import js from '@eslint/js';
import tsParser from '@typescript-eslint/parser';
import tsPlugin from '@typescript-eslint/eslint-plugin';
import svelte from 'eslint-plugin-svelte';

const browserGlobals = {
	console: 'readonly',
	document: 'readonly',
	Event: 'readonly',
	SubmitEvent: 'readonly',
	window: 'readonly'
};

export default [
	js.configs.recommended,
	...svelte.configs['flat/recommended'],
	{
		ignores: ['.svelte-kit/**', 'build/**', 'node_modules/**']
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
			...tsPlugin.configs.recommended.rules
		}
	},
	{
		files: ['**/*.svelte'],
		languageOptions: {
			parserOptions: {
				parser: tsParser
			},
			globals: browserGlobals
		}
	}
];
