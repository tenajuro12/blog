// e2e/helpers.js
// Shared helpers for all E2E tests

const BASE_URL = 'http://localhost:3000';
const API_URL  = 'http://localhost:8080';

/**
 * Register a new user via API (faster than going through UI each time)
 * Returns { token, user }
 */
async function registerUser(request, suffix = Date.now()) {
  const res = await request.post(`${API_URL}/api/auth/register`, {
    data: {
      username: `testuser_${suffix}`,
      email:    `test_${suffix}@example.com`,
      password: 'password123',
    },
  });
  const body = await res.json();
  return {
    token:    body.token,
    user:     body.user,
    email:    `test_${suffix}@example.com`,
    password: 'password123',
    username: `testuser_${suffix}`,
  };
}

/**
 * Log in through the UI — fills email + password, clicks Sign in
 * Returns after redirect to home page
 */
async function loginUI(page, email, password) {
  await page.goto('/login');
  // inputs have no id/name — select by placeholder
  await page.getByPlaceholder('you@example.com').fill(email);
  await page.getByPlaceholder('••••••••').fill(password);
  await page.getByRole('button', { name: /sign in/i }).click();
  // wait for redirect away from /login
  await page.waitForURL(url => !url.pathname.includes('/login'), { timeout: 8000 });
}

/**
 * Create a post via API
 * Returns the created post object (including slug)
 */
async function createPost(request, token, title, body, tags = '') {
  const res = await request.post(`${API_URL}/api/posts`, {
    headers: { Authorization: `Bearer ${token}` },
    data: { title, body, tags },
  });
  return await res.json();
}

module.exports = { BASE_URL, API_URL, registerUser, loginUI, createPost };
