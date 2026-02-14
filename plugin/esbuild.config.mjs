import { build, context } from 'esbuild';
import { readFileSync, writeFileSync, mkdirSync } from 'fs';

const isWatch = process.argv.includes('--watch');

// Build plugin code (runs in Figma sandbox)
const codeConfig = {
  entryPoints: ['src/main.ts'],
  bundle: true,
  outfile: 'dist/code.js',
  target: 'es6',
  format: 'iife',
  minify: !isWatch,
};

// Build UI code
const uiConfig = {
  entryPoints: ['src/ui/app.ts'],
  bundle: true,
  outfile: 'dist/ui.js',
  target: 'es6',
  format: 'iife',
  minify: !isWatch,
};

async function buildAll() {
  mkdirSync('dist', { recursive: true });

  await build(codeConfig);
  await build(uiConfig);

  // Read CSS and JS, inline into HTML
  const css = readFileSync('src/ui/styles.css', 'utf8');
  const js = readFileSync('dist/ui.js', 'utf8');
  // Use function callbacks to avoid $& back-reference interpretation in minified JS
  const html = readFileSync('src/ui/index.html', 'utf8')
    .replace('/* STYLES */', () => css)
    .replace('/* SCRIPT */', () => js);

  writeFileSync('dist/ui.html', html);
  console.log('Build complete');
}

if (isWatch) {
  const codeCtx = await context(codeConfig);
  const uiCtx = await context(uiConfig);
  await codeCtx.watch();
  await uiCtx.watch();

  // Do an initial full build for the HTML inlining
  await buildAll();
  console.log('Watching for changes...');
} else {
  await buildAll();
}
