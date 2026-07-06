import { defineConfig } from 'vitest/config';
import viteConfig from './vite.config';

export default defineConfig({
	...viteConfig,
	test: {
		include: ['src/**/*.test.ts'],
		coverage: {
			provider: 'v8',
			reportsDirectory: 'coverage',
			reporter: ['text', 'html', 'lcov'],
			include: ['src/**/*.{ts,svelte}'],
			exclude: [
				'src/lib/api/generated/**',
				'src/**/*.d.ts',
				'src/routes/**',
				'src/lib/features/**',
				'src/lib/components/ui/**',
				'src/lib/components/**/*Modal.svelte',
				'src/lib/components/**/*Picker.svelte',
				'src/lib/components/**/*Sheet.svelte',
				'src/lib/components/**/*Select.svelte',
				'src/lib/components/**/*Autocomplete.svelte'
			],
			thresholds: {
				statements: 60
			}
		}
	}
});
