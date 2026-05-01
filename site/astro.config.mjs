import { defineConfig } from 'astro/config';
import starlight from '@astrojs/starlight';

export default defineConfig({
  site: 'https://aihappydesign.com',
  integrations: [
    starlight({
      title: 'AI Happy Design',
      description: 'Agent-ready Figma automation with 184 schema-backed commands, MCP support, token export, accessibility audits, and design-code parity checks.',
      favicon: '/favicon.svg',
      customCss: ['./src/styles/site.css'],
      social: [
        { icon: 'github', label: 'GitHub', href: 'https://github.com/nerveband/ai-happy-design' },
        { icon: 'twitter', label: 'X / Twitter', href: 'https://x.com/workandpromise' }
      ],
      sidebar: [
        {
          label: 'Getting Started',
          items: [
            { label: 'Installation', slug: 'getting-started/installation' },
            { label: 'Quick Start', slug: 'getting-started/quickstart' }
          ]
        },
        {
          label: 'Guides',
          items: [
            { label: 'CLI Usage', slug: 'guides/cli-usage' },
            { label: 'Batch Operations', slug: 'guides/batch-operations' },
            { label: 'Design Intelligence', slug: 'guides/design-intelligence' },
            { label: 'Release v0.13.2', slug: 'guides/release-v0132' }
          ]
        },
        {
          label: 'Reference',
          items: [
            { label: 'Command Reference', slug: 'reference/commands' },
            { label: 'Schema System', slug: 'reference/schema' },
            { label: 'LLM Files', slug: 'reference/llm-files' }
          ]
        },
        {
          label: 'Use Cases',
          items: [
            { label: 'Design Systems', slug: 'use-cases/design-systems' },
            { label: 'Marketing', slug: 'use-cases/marketing' },
            { label: 'Multi-Agent Workflows', slug: 'use-cases/multi-agent' }
          ]
        },
        {
          label: 'Architecture',
          items: [
            { label: 'Overview', slug: 'architecture/overview' }
          ]
        }
      ],
      tableOfContents: { minHeadingLevel: 2, maxHeadingLevel: 3 },
      editLink: {
        baseUrl: 'https://github.com/nerveband/ai-happy-design/edit/main/site/'
      },
      head: [
        {
          tag: 'script',
          attrs: {
            defer: true,
            src: 'https://stats.wavedepth.com/script.js',
            'data-website-id': 'fe025f86-6bfb-4fa1-9d4a-eb344b55037b'
          }
        }
      ]
    })
  ]
});
