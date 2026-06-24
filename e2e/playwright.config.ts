import { defineConfig, devices } from '@playwright/test';
import { fileURLToPath } from 'node:url';
import path from 'node:path';

// Resolve repo paths relative to this config file so the webServer commands work
// regardless of the directory Playwright is invoked from.
const __dirname = path.dirname(fileURLToPath(import.meta.url));
const repoRoot = path.resolve(__dirname, '..');

const BACKEND_PORT = 8080;
const FRONTEND_PORT = 5173;
const BACKEND_URL = `http://localhost:${BACKEND_PORT}`;
const FRONTEND_URL = `http://localhost:${FRONTEND_PORT}`;

/**
 * Integration tests run the real app: the Go backend (seeded, in-memory) and the
 * Vue SPA, wired together so the SPA talks to the backend over the HTTP contract.
 * Playwright boots both via `webServer` and drives a real browser against the SPA.
 * See docs/adr/0003-integration-tests-playwright.md.
 */
export default defineConfig({
  testDir: './tests',
  // Shared in-memory backend state means tests must run serially and predictably.
  fullyParallel: false,
  workers: 1,
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 1 : 0,
  reporter: [['list'], ['html', { open: 'never' }]],

  use: {
    baseURL: FRONTEND_URL,
    // The API speaks UTC; render times in UTC so assertions are timezone-stable.
    timezoneId: 'UTC',
    trace: 'on-first-retry',
    screenshot: 'only-on-failure',
  },

  projects: [
    {
      name: 'chromium',
      use: { ...devices['Desktop Chrome'] },
    },
  ],

  webServer: [
    {
      // Go backend with its built-in demo seed (event types + availability).
      command: 'go run ./cmd/server',
      cwd: path.join(repoRoot, 'backend'),
      url: `${BACKEND_URL}/event-types`,
      env: {
        PORT: String(BACKEND_PORT),
        SEED_DATA: 'true',
        CORS_ALLOWED_ORIGIN: '*',
      },
      reuseExistingServer: !process.env.CI,
      timeout: 180_000,
      stdout: 'pipe',
      stderr: 'pipe',
    },
    {
      // Vue SPA pointed at the local backend instead of the Prism mock default.
      // Vite only exposes VITE_* from .env files (not arbitrary process env), so
      // `--mode e2e` loads frontend/.env.e2e (VITE_API_BASE_URL=:8080).
      command: 'npm run dev -- --port 5173 --strictPort --mode e2e',
      cwd: path.join(repoRoot, 'frontend'),
      url: FRONTEND_URL,
      reuseExistingServer: !process.env.CI,
      timeout: 120_000,
      stdout: 'pipe',
      stderr: 'pipe',
    },
  ],
});
