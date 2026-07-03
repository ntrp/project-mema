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
			exclude: ['src/lib/api/generated/**', 'src/**/*.d.ts'],
			thresholds: {
				statements: 60
			}
		}
	}
});
