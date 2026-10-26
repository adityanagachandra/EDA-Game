import { test, expect } from './fixtures'

// Regression: layout + page both rendered chrome on these routes
// (double headers on /docs, double bottom tab bars on /diagnose + /goals).
// Pages own their chrome; layouts are metadata-only (docs/layout keeps its
// section header by design).
test.beforeEach(async ({ context }) => {
  await context.route('**/api/**', route => route.fulfill({ status: 404, body: '' }))
})

for (const route of ['/docs', '/how-it-works']) {
  test(`${route} renders exactly one public header`, async ({ page }) => {
    await page.goto(route)
    await expect(page.locator('header')).toHaveCount(1, { timeout: 30_000 })
    await expect(page.getByTestId('profile-sidebar')).toHaveCount(0)
    await expect(page.locator('footer')).toHaveCount(1)
    if (route === '/docs') {
      // Docs section header survives; the default app header must not.
      await expect(page.getByText('Architecture, design and the contributing guide')).toBeVisible()
    }
  })
}

for (const route of ['/diagnose', '/goals']) {
  test(`${route} renders exactly one app shell and no public header`, async ({ page }) => {
    await page.goto(route)
    await expect(page.locator('header')).toHaveCount(0, { timeout: 30_000 })
    await expect(page.getByTestId('profile-sidebar')).toHaveCount(1)
    await expect(page.locator('footer')).toHaveCount(1)
    await expect(page.locator('nav.lg\\:hidden')).toHaveCount(1)
  })
}
