import { test, expect } from './fixtures'

const USER = {
  student_id: 's1',
  name: 'Tester',
  username: 'tester',
  concepts_mastered: 3,
  current_streak: 2,
  level: 'beginner',
  diagnostic_completed: true,
}

test.beforeEach(async ({ context }) => {
  await context.route('**/api/**', route => route.fulfill({ status: 404, body: '' }))
})

test('header shows Login for anonymous visitors', async ({ page }) => {
  // /privacy uses the default public Header. /how-it-works passes a custom
  // link list without Login, and product pages have no Header at all.
  await page.goto('/privacy')
  await expect(page.locator('header').getByText('Login').first()).toBeVisible()
})

test('header shows authenticated links when a token exists', async ({ page }) => {
  await page.addInitScript(([user]) => {
    localStorage.setItem('mathua_token', 'fake-token')
    localStorage.setItem('mathua_user', JSON.stringify(user))
  }, [USER])
  await page.goto('/privacy')
  const header = page.locator('header')
  await expect(header.getByText('Profile').first()).toBeVisible()
  await expect(header.getByText('Learn').first()).toBeVisible()
  await expect(header.getByText('Login')).toHaveCount(0)
})
