import { mkdirSync, copyFileSync } from 'node:fs';

mkdirSync('static/data', { recursive: true });
copyFileSync(
  '../scripts/pattern_descriptions/pattern_descriptions.json',
  'static/data/pattern_descriptions.json'
);
