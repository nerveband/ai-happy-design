// @ts-check
import { defineConfig } from 'astro/config';
import starlight from '@astrojs/starlight';
// beautiful-mermaid is used directly in .astro components via renderMermaid()

// https://astro.build/config
export default defineConfig({
	site: 'https://aihappydesign.com',
	integrations: [
		starlight({
			title: 'AI Happy Design',
			tagline: 'Give LLMs full Figma canvas access',
			logo: {
				light: './src/assets/logo-light.svg',
				dark: './src/assets/logo-dark.svg',
			},
			social: [
				{
					icon: 'github',
					label: 'GitHub',
					href: 'https://github.com/nerveband/ai-happy-design',
				},
				{
					icon: 'x.com',
					label: 'X/Twitter',
					href: 'https://x.com/workandpromise',
				},
			],
			customCss: [
				'./src/styles/custom.css',
			],
			head: [
				{
					tag: 'link',
					attrs: {
						rel: 'preconnect',
						href: 'https://fonts.googleapis.com',
					},
				},
				{
					tag: 'link',
					attrs: {
						rel: 'preconnect',
						href: 'https://fonts.gstatic.com',
						crossorigin: 'anonymous',
					},
				},
				{
					tag: 'link',
					attrs: {
						href: 'https://fonts.googleapis.com/css2?family=Space+Grotesk:wght@400;500;600;700&family=Inter:wght@400;500;600;700&family=JetBrains+Mono:wght@400;500;700&display=swap',
						rel: 'stylesheet',
					},
				},
				{
					tag: 'link',
					attrs: {
						rel: 'icon',
						type: 'image/svg+xml',
						href: "data:image/svg+xml,<svg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 28 28'><rect width='28' height='28' rx='6' fill='%23FFD600'/><path d='M15.5 4L8 15h5l-1.5 9L20 13h-5l1.5-9z' fill='%230a0a0a'/></svg>",
					},
				},
				{
					tag: 'script',
					attrs: {
						defer: true,
						src: 'https://stats.wavedepth.com/script.js',
						'data-website-id': 'fe025f86-6bfb-4fa1-9d4a-eb344b55037b',
					},
				},
			],
			sidebar: [
				{
					label: 'Getting Started',
					items: [
						{ label: 'Installation', slug: 'getting-started/installation' },
						{ label: 'Quick Start', slug: 'getting-started/quickstart' },
					],
				},
				{
					label: 'Guides',
					items: [
						{ label: 'CLI Usage', slug: 'guides/cli-usage' },
						{ label: 'Batch Operations', slug: 'guides/batch-operations' },
						{ label: 'Design Intelligence', slug: 'guides/design-intelligence' },
						{ label: 'Design Patterns', slug: 'guides/design-patterns' },
					],
				},
				{
					label: 'Reference',
					items: [
						{ label: 'Command Reference', slug: 'reference/commands' },
						{ label: 'Schema System', slug: 'reference/schema' },
						{ label: 'Batch Aliases', slug: 'reference/batch-aliases' },
					],
				},
				{
					label: 'Use Cases',
					items: [
						{ label: 'Social Media', slug: 'use-cases/social-media' },
						{ label: 'Presentations', slug: 'use-cases/presentations' },
						{ label: 'Design Systems', slug: 'use-cases/design-systems' },
						{ label: 'Marketing', slug: 'use-cases/marketing' },
						{ label: 'Prototyping', slug: 'use-cases/prototyping' },
						{ label: 'Multi-Agent Workflows', slug: 'use-cases/multi-agent' },
					],
				},
				{
					label: 'Architecture',
					items: [
						{ label: 'Overview', slug: 'architecture/overview' },
						{ label: 'Protocol', slug: 'architecture/protocol' },
					],
				},
			],
		}),
	],
});
