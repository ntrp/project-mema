import { defineConfig, devices } from '@playwright/test';

const baseURL = process.env.E2E_BASE_URL ?? 'http://127.0.0.1:15173';
const reuseExistingServer = !process.env.CI;

export default defineConfig({
	testDir: 'e2e',
	timeout: 30_000,
	expect: { timeout: 10_000 },
	use: {
		baseURL,
		trace: 'retain-on-failure'
	},
	webServer: [
		{
			command: 'ADDR=0.0.0.0:18080 make dev-api',
			cwd: '..',
			url: 'http://127.0.0.1:18080/api/health',
			reuseExistingServer,
			timeout: 120_000
		},
		{
			command: 'node e2e/support/mock-services.mjs',
			url: 'http://127.0.0.1:18180/health',
			reuseExistingServer,
			timeout: 30_000
		},
		{
			command:
				'NVIM_LISTEN_ADDRESS=/tmp/project-mema.nvim LAUNCH_EDITOR=/Users/ntrp/_pws/project-mema/scripts/open-in-nvim.sh pnpm exec vite dev --host 0.0.0.0 --port 15173',
			url: baseURL,
			reuseExistingServer,
			timeout: 120_000
		}
	],
	projects: [
		{
			name: 'chromium',
			use: { ...devices['Desktop Chrome'] }
		}
	]
});
