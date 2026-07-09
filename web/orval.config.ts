import { defineConfig } from 'orval';

export default defineConfig({
	mediaManager: {
		input: {
			target: '../api/openapi.yaml',
			unsafeDisableValidation: true
		},
		output: {
			baseUrl: '/api',
			client: 'svelte-query',
			httpClient: 'fetch',
			mode: 'single',
			override: {
				fetch: {
					includeHttpResponseReturnType: false
				},
				query: {
					useQuery: true
				}
			},
			target: 'src/lib/api/generated/tanstack.ts'
		}
	}
});
