import { playwright } from '@vitest/browser-playwright';
import { defineConfig } from 'vitest/config';
import viteConfig from './vite.config';

export default defineConfig({
	...viteConfig,
	test: {
		include: ['src/**/*.browser.test.ts'],
		browser: {
			enabled: true,
			headless: true,
			provider: playwright(),
			instances: [{ browser: 'chromium' }]
		},
		coverage: {
			provider: 'v8',
			reportsDirectory: 'coverage/browser',
			reporter: ['text', 'html', 'lcov'],
			include: [
				'src/lib/components/app/renderer/PlayOnDevice.svelte',
				'src/lib/components/app/renderer/RendererRemote.svelte'
			],
			thresholds: {
				perFile: true,
				statements: 60,
				branches: 60,
				functions: 60,
				lines: 60
			}
		}
	}
});
