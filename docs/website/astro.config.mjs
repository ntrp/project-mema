import { defineConfig } from 'astro/config';
import starlight from '@astrojs/starlight';
import starlightSidebarTopics from 'starlight-sidebar-topics';

export default defineConfig({
	integrations: [
		starlight({
			title: 'Media Manager Docs',
			logo: {
				src: './src/assets/icon.svg',
				alt: 'Media Manager'
			},
			social: [{ icon: 'github', label: 'GitHub', href: 'https://github.com/ntrp/project-mema' }],
			plugins: [
				starlightSidebarTopics([
					{
						id: 'user-guide',
						label: 'User Guide',
						link: '/user-guide/getting-started/first-run/',
						icon: 'open-book',
						items: [
							{
								label: 'Getting Started',
								items: [
									{ label: 'First Run', link: '/user-guide/getting-started/first-run/' },
									{
										label: 'Configuration',
										link: '/user-guide/getting-started/configuration/'
									}
								]
							},
							{
								label: 'Concepts',
								items: [
									{ label: 'Media Lifecycle', link: '/user-guide/concepts/media-lifecycle/' },
									{ label: 'Profiles', link: '/user-guide/concepts/profiles/' },
									{ label: 'Library Import', link: '/user-guide/concepts/library-import/' },
									{ label: 'Track Management', link: '/user-guide/concepts/track-management/' }
								]
							}
						]
					},
					{
						id: 'dev-guide',
						label: 'Dev Guide',
						link: '/dev-guide/development-workflow/',
						icon: 'setting',
						items: [
							{
								label: 'Development',
								items: [
									{
										label: 'Development Workflow',
										link: '/dev-guide/development-workflow/'
									}
								]
							}
						]
					},
					{
						id: 'architecture',
						label: 'Architecture',
						link: '/architecture/system-overview/',
						icon: 'information',
						items: [
							{
								label: 'Architecture',
								items: [
									{ label: 'System Overview', link: '/architecture/system-overview/' },
									{ label: 'Storage', link: '/architecture/storage/' },
									{ label: 'API Contract', link: '/architecture/api-contract/' }
								]
							}
						]
					}
				])
			],
			customCss: ['./src/styles/custom.css']
		})
	]
});
