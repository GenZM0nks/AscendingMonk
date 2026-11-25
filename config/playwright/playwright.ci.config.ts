import { defineConfig } from '@playwright/test';
import baseConfig from './playwright.config';

export default defineConfig({
    ...baseConfig, // Uses playwright.config.ts
    testDir: '../../test/e2e',
});