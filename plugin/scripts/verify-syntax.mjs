import { readFileSync } from 'fs';

const codePath = 'dist/code.js';
const code = readFileSync(codePath, 'utf8');

const checks = [
  {
    name: 'optional chaining',
    regex: /\?\./g,
    guidance: 'Replace optional chaining (?.) with explicit guards.',
  },
  {
    name: 'nullish coalescing',
    regex: /\?\?/g,
    guidance: 'Replace ?? with explicit null/undefined checks.',
  },
  {
    name: 'spread syntax',
    regex: /\.\.\./g,
    guidance: 'Avoid spread syntax in plugin runtime code.',
  },
  {
    name: 'for-await loops',
    regex: /for\s+await\s*\(/g,
    guidance: 'Replace for-await loops with Promise-based iteration.',
  },
  {
    name: 'logical assignment',
    regex: /(\?\?=|\|\|=|&&=)/g,
    guidance: 'Replace logical assignments with explicit assignments.',
  },
  {
    name: 'private class fields',
    regex: /#[A-Za-z_$][\w$]*/g,
    guidance: 'Avoid private class fields in plugin runtime code.',
  },
];

const violations = [];

for (const check of checks) {
  const matches = code.match(check.regex);
  if (matches && matches.length > 0) {
    violations.push(`${check.name}: ${matches.length} occurrence(s). ${check.guidance}`);
  }
}

if (violations.length > 0) {
  console.error('Figma QuickJS compatibility check failed for ' + codePath);
  for (const violation of violations) {
    console.error(' - ' + violation);
  }
  process.exit(1);
}

console.log('Syntax compatibility check passed for ' + codePath);
