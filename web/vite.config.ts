import adapter from '@sveltejs/adapter-static';
import { sveltekit } from '@sveltejs/kit/vite';
import tailwindcss from '@tailwindcss/vite';
import { env } from 'node:process';
import { defineConfig } from 'vite';

const apiProxyTarget = env.VITE_API_PROXY_TARGET ?? 'http://127.0.0.1:18080';

export default defineConfig({
  plugins: [
    tailwindcss(),
    sveltekit({
      compilerOptions: {
        // Force runes mode for the project, except for libraries. Can be removed in svelte 6.
        runes: ({ filename }) =>
          filename.split(/[/\\]/).includes('node_modules') ? undefined : true
      },
      inspector: {
        showToggleButton: 'always',
        toggleButtonPos: 'bottom-right'
      },
      adapter: adapter({
        fallback: '200.html'
      })
    })
  ],
  server: {
    port: 15173,
    strictPort: true,
    proxy: {
      '/api': apiProxyTarget
    }
  }
});
