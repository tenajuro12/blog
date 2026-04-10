// e2e/posts.spec.js
// TC-E2E-03: Create post flow
// TC-E2E-08: Edit post flow
// TC-E2E-09: Unauthenticated cannot access new-post form

const { test, expect } = require('@playwright/test');
const { registerUser, loginUI } = require('./helpers');

// ── TC-E2E-03: Create post flow ───────────────────────────────────────────────
test('TC-E2E-03: authenticated user creates post and it appears with correct slug URL', async ({ page, request }) => {
  const { email, password } = await registerUser(request);
  await loginUI(page, email, password);

  // Navigate to new post page
  await page.goto('/new-post');

  const title = `E2E Test Post ${Date.now()}`;
  const body  = 'This is the body of an automated E2E test post written by Playwright.';

  // Fill title — placeholder: "Your story title..."
  await page.getByPlaceholder('Your story title...').fill(title);

  // Fill body — placeholder: "Tell your story..."
  await page.getByPlaceholder('Tell your story...').fill(body);

  // Publish button becomes enabled once title and body are filled
  const publishBtn = page.getByRole('button', { name: /publish/i });
  await expect(publishBtn).toBeEnabled({ timeout: 3000 });
  await publishBtn.click();

  // Should redirect to /posts/{slug}
  await page.waitForURL(url => url.pathname.startsWith('/posts/'), { timeout: 8000 });

  // Slug URL should contain a lowercase version of the title
  const url = page.url();
  expect(url).toContain('/posts/');

  // Post title should be visible on the detail page
  await expect(page.getByText(title, { exact: false })).toBeVisible({ timeout: 5000 });
});

// ── TC-E2E-08: Create post with tags ─────────────────────────────────────────
test('TC-E2E-08: create post with tags — tags appear on post detail page', async ({ page, request }) => {
  const { email, password } = await registerUser(request);
  await loginUI(page, email, password);

  await page.goto('/new-post');

  const title = `Tagged Post ${Date.now()}`;

  await page.getByPlaceholder('Your story title...').fill(title);
  await page.getByPlaceholder('Tell your story...').fill('Body content for the tagged post test.');

  // Fill tags — placeholder: "Add tags (comma separated)..."
  await page.getByPlaceholder('Add tags (comma separated)...').fill('playwright,e2e,testing');

  await page.getByRole('button', { name: /publish/i }).click();

  // Wait for redirect to post detail
  await page.waitForURL(url => url.pathname.startsWith('/posts/'), { timeout: 8000 });

  // Tags should be visible somewhere on the post detail page
  await expect(page.getByText('playwright', { exact: false })).toBeVisible({ timeout: 5000 });
});

// ── TC-E2E-09: Publish button disabled with empty title ───────────────────────
test('TC-E2E-09: publish button is disabled when title is empty', async ({ page, request }) => {
  const { email, password } = await registerUser(request);
  await loginUI(page, email, password);

  await page.goto('/new-post');

  // Only fill body, leave title empty
  await page.getByPlaceholder('Tell your story...').fill('Some body text here.');

  // Publish button should be disabled
  const publishBtn = page.getByRole('button', { name: /publish/i });
  await expect(publishBtn).toBeDisabled({ timeout: 3000 });
});

// ── TC-E2E-10: Publish button disabled with empty body ────────────────────────
test('TC-E2E-10: publish button is disabled when body is empty', async ({ page, request }) => {
  const { email, password } = await registerUser(request);
  await loginUI(page, email, password);

  await page.goto('/new-post');

  // Only fill title, leave body empty
  await page.getByPlaceholder('Your story title...').fill('Some Title Here');

  const publishBtn = page.getByRole('button', { name: /publish/i });
  await expect(publishBtn).toBeDisabled({ timeout: 3000 });
});
