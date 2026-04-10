// e2e/auth.spec.js
// TC-E2E-01: Registration flow
// TC-E2E-02: Login flow
// TC-E2E-05: Login failure flow
// TC-E2E-06: Protected route redirect

const { test, expect } = require('@playwright/test');
const { registerUser, loginUI } = require('./helpers');

// ── TC-E2E-01: Registration flow ──────────────────────────────────────────────
test('TC-E2E-01: successful registration redirects to home and shows username', async ({ page, request }) => {
  const suffix = Date.now();
  const username = `newuser_${suffix}`;
  const email    = `newuser_${suffix}@example.com`;
  const password = 'password123';

  await page.goto('/register');

  // Fill username, email, password by placeholder
  await page.getByPlaceholder('yourname').fill(username);
  await page.getByPlaceholder('you@example.com').fill(email);
  await page.getByPlaceholder('••••••••').fill(password);

  // Click Create account button
  await page.getByRole('button', { name: /create account/i }).click();

  // Should redirect away from /register
  await page.waitForURL(url => !url.pathname.includes('/register'), { timeout: 8000 });

  // Should land on home page
  expect(page.url()).toContain('localhost:3000');
  expect(page.url()).not.toContain('/register');

  // Username should appear somewhere in the nav
  await expect(page.getByText(username, { exact: false })).toBeVisible({ timeout: 5000 });
});

// ── TC-E2E-02: Login flow ─────────────────────────────────────────────────────
test('TC-E2E-02: successful login redirects to home and stores token', async ({ page, request }) => {
  // Create user via API so we have known credentials
  const { email, password, username } = await registerUser(request);

  await page.goto('/login');
  await page.getByPlaceholder('you@example.com').fill(email);
  await page.getByPlaceholder('••••••••').fill(password);
  await page.getByRole('button', { name: /sign in/i }).click();

  // Wait for redirect
  await page.waitForURL(url => !url.pathname.includes('/login'), { timeout: 8000 });

  // Landed on home
  expect(page.url()).not.toContain('/login');

  // Token stored in localStorage
  const token = await page.evaluate(() => localStorage.getItem('token'));
  expect(token).toBeTruthy();
  expect(typeof token).toBe('string');
  expect(token.length).toBeGreaterThan(20);

  // Username visible in UI
  await expect(page.getByText(username, { exact: false })).toBeVisible({ timeout: 5000 });
});

// ── TC-E2E-05: Login failure flow ─────────────────────────────────────────────
test('TC-E2E-05: wrong password shows error and stays on /login', async ({ page, request }) => {
  const { email } = await registerUser(request);

  await page.goto('/login');
  await page.getByPlaceholder('you@example.com').fill(email);
  await page.getByPlaceholder('••••••••').fill('wrongpassword999');
  await page.getByRole('button', { name: /sign in/i }).click();

  // Should NOT redirect — still on /login
  await page.waitForTimeout(1500);
  expect(page.url()).toContain('/login');

  // Error message should be visible
  // The Auth.js component renders error in a styled div
  await expect(page.getByText(/invalid credentials/i).or(page.getByText(/incorrect/i)).or(page.getByText(/unauthorized/i)).or(page.getByText(/error/i).first())).toBeVisible({ timeout: 5000 });

  // Token should NOT be stored
  const token = await page.evaluate(() => localStorage.getItem('token'));
  expect(token).toBeFalsy();
});

// ── TC-E2E-06: Protected route redirect ───────────────────────────────────────
test('TC-E2E-06: unauthenticated user navigating to /new-post is redirected to /login', async ({ page }) => {
  // Start with clean state — no token
  await page.goto('/');
  await page.evaluate(() => localStorage.removeItem('token'));
  await page.evaluate(() => localStorage.removeItem('user'));

  // Try to access protected route directly
  await page.goto('/new-post');

  // Should be redirected to /login
  await page.waitForURL(url => url.pathname.includes('/login'), { timeout: 8000 });
  expect(page.url()).toContain('/login');
});

// ── TC-E2E-07: Register with missing fields shows validation ──────────────────
test('TC-E2E-07: register with empty email stays on register page', async ({ page }) => {
  await page.goto('/register');

  // Fill only username and password, leave email empty
  await page.getByPlaceholder('yourname').fill('someuser');
  await page.getByPlaceholder('••••••••').fill('password123');

  await page.getByRole('button', { name: /create account/i }).click();

  // Should stay on /register (either browser validation or API error)
  await page.waitForTimeout(1500);
  expect(page.url()).toContain('/register');
});
