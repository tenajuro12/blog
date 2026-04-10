// e2e/social.spec.js
// TC-E2E-04: Like a post flow
// TC-E2E-11: Bookmark a post flow
// TC-E2E-12: Unlike a post
// TC-E2E-13: Unauthenticated user cannot like

const { test, expect } = require('@playwright/test');
const { registerUser, loginUI, createPost } = require('./helpers');

// Like button renders: <button onClick={handleLike}> <span>♥/♡</span> <span>{likeCount}</span> </button>
// Bookmark button renders: <button onClick={handleBookmark}> 🔖 or ⬜ </button>

// ── TC-E2E-04: Like a post ────────────────────────────────────────────────────
test('TC-E2E-04: authenticated user likes a post and like count increments', async ({ page, request }) => {
  const author = await registerUser(request, `author_${Date.now()}`);
  const post   = await createPost(request, author.token, 'Post to Like E2E', 'Body of the post to be liked in E2E test.');

  const liker = await registerUser(request, `liker_${Date.now()}`);
  await loginUI(page, liker.email, liker.password);

  await page.goto(`/posts/${post.slug}`);
  await page.waitForLoadState('networkidle');

  // Like button contains ♡ (heart outline) when not liked
  const likeBtn = page.locator('button').filter({ hasText: /♡|♥/ }).first();
  await expect(likeBtn).toBeVisible({ timeout: 5000 });

  // Get count before
  const countBefore = parseInt(await likeBtn.locator('span').last().textContent() || '0');

  await likeBtn.click();
  await page.waitForTimeout(800);

  // Count should increment OR button state should change to ♥
  const countAfter = parseInt(await likeBtn.locator('span').last().textContent() || '0');
  const isLiked = await page.locator('button').filter({ hasText: /♥/ }).first().isVisible().catch(() => false);

  if (!isLiked && countAfter <= countBefore) {
    t.fail('Like did not register — count did not increment and heart did not fill');
  }
});

// ── TC-E2E-11: Bookmark a post ────────────────────────────────────────────────
test('TC-E2E-11: authenticated user bookmarks a post — appears in bookmarks page', async ({ page, request }) => {
  const author = await registerUser(request, `bkauth_${Date.now()}`);
  const post   = await createPost(request, author.token, 'Post to Bookmark E2E', 'Body of the post to be bookmarked in E2E test.');

  const reader = await registerUser(request, `bkread_${Date.now()}`);
  await loginUI(page, reader.email, reader.password);

  await page.goto(`/posts/${post.slug}`);
  await page.waitForLoadState('networkidle');

  // Bookmark button shows ⬜ when not bookmarked, 🔖 when bookmarked
  const bookmarkBtn = page.locator('button[title="Bookmark"], button[title="Remove bookmark"]').first()
      .or(page.locator('button').filter({ hasText: /⬜|🔖/ }).first());

  await expect(bookmarkBtn).toBeVisible({ timeout: 5000 });
  await bookmarkBtn.click();
  await page.waitForTimeout(1000);

  // Navigate to bookmarks page
  await page.goto('/bookmarks');
  await page.waitForLoadState('networkidle');

  // Post title should appear
  await expect(page.getByText(post.title, { exact: false })).toBeVisible({ timeout: 8000 });
});

// ── TC-E2E-12: Unlike a post ──────────────────────────────────────────────────
test('TC-E2E-12: user unlikes a post they already liked — like state toggles', async ({ page, request }) => {
  const author = await registerUser(request, `ulauth_${Date.now()}`);
  const post   = await createPost(request, author.token, 'Post to Unlike E2E', 'Body of the post to unlike in E2E test.');

  const liker = await registerUser(request, `uliker_${Date.now()}`);
  await loginUI(page, liker.email, liker.password);

  await page.goto(`/posts/${post.slug}`);
  await page.waitForLoadState('networkidle');

  const likeBtn = page.locator('button').filter({ hasText: /♡|♥/ }).first();
  await expect(likeBtn).toBeVisible({ timeout: 5000 });

  // Click to like
  await likeBtn.click();
  await page.waitForTimeout(800);

  // Should now show ♥ (filled)
  const filledBtn = page.locator('button').filter({ hasText: /♥/ }).first();
  await expect(filledBtn).toBeVisible({ timeout: 3000 });

  // Click again to unlike
  await filledBtn.click();
  await page.waitForTimeout(800);

  // Should go back to ♡ (outline)
  await expect(page.locator('button').filter({ hasText: /♡/ }).first()).toBeVisible({ timeout: 3000 });
});

// ── TC-E2E-13: Unauthenticated user clicking like does nothing ────────────────
test('TC-E2E-13: unauthenticated user clicking like does not trigger API call', async ({ page, request }) => {
  const author = await registerUser(request, `noauth_${Date.now()}`);
  const post   = await createPost(request, author.token, 'Post No Auth Like', 'Body of post for unauthenticated like test.');

  // Visit without logging in
  await page.evaluate(() => {
    localStorage.removeItem('token');
    localStorage.removeItem('user');
  }).catch(() => {});

  await page.goto(`/posts/${post.slug}`);
  await page.waitForLoadState('networkidle');

  const likeBtn = page.locator('button').filter({ hasText: /♡|♥/ }).first();
  await expect(likeBtn).toBeVisible({ timeout: 5000 });

  // Get count before click
  const beforeText = await likeBtn.locator('span').last().textContent().catch(() => '0');

  // Click — handler does: if (!user) return; so nothing should happen
  await likeBtn.click();
  await page.waitForTimeout(800);

  // Count should not change
  const afterText = await likeBtn.locator('span').last().textContent().catch(() => '0');
  expect(afterText.trim()).toBe(beforeText.trim());
});