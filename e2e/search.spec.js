const { test, expect } = require("./fixtures");

test.describe("Search", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto("/");
  });

  test("renders the search button in topnav", async ({ page }) => {
    const btn = page.locator('nav[data-topnav] button[aria-label="Search"]');
    await expect(btn).toBeVisible();
    await expect(btn).toContainText("Search");
    await expect(page.locator("nav[data-topnav] kbd")).toContainText("/");
  });

  test("opens search dialog on button click", async ({ page }) => {
    await expect(page.locator("#search-dialog")).not.toBeVisible();
    await page.locator('button[aria-label="Search"]').click();
    await expect(page.locator("#search-dialog")).toBeVisible();
    await expect(
      page.locator("#search-dialog input[type='search']")
    ).toBeFocused();
  });

  test("opens search dialog with / shortcut", async ({ page }) => {
    await page.keyboard.press("/");
    await expect(page.locator("#search-dialog")).toBeVisible();
    await expect(
      page.locator("#search-dialog input[type='search']")
    ).toBeFocused();
  });

  test("closes search dialog with Escape", async ({ page }) => {
    await page.locator('button[aria-label="Search"]').click();
    await expect(page.locator("#search-dialog")).toBeVisible();
    await page.keyboard.press("Escape");
    await expect(page.locator("#search-dialog")).not.toBeVisible();
  });

  test("shows status message before typing", async ({ page }) => {
    await page.locator('button[aria-label="Search"]').click();
    await expect(page.locator("[data-search-status]")).toContainText(
      "Type at least 2 characters"
    );
  });

  test("returns results for a valid query", async ({ page }) => {
    await page.locator('button[aria-label="Search"]').click();
    await page.locator("#search-dialog input[type='search']").fill("config");
    // Wait for results to load
    await expect(page.locator(".search-result").first()).toBeVisible({
      timeout: 5000,
    });
    const count = await page.locator(".search-result").count();
    expect(count).toBeGreaterThanOrEqual(1);
  });

  test("shows result titles and summaries", async ({ page }) => {
    await page.locator('button[aria-label="Search"]').click();
    await page.locator("#search-dialog input[type='search']").fill("config");
    await expect(page.locator(".search-result").first()).toBeVisible({
      timeout: 5000,
    });
    const firstResult = page.locator(".search-result").first();
    await expect(firstResult.locator("strong")).toBeAttached();
    await expect(firstResult.locator("small")).toBeAttached();
  });

  test("navigates to a result on click", async ({ page }) => {
    await page.locator('button[aria-label="Search"]').click();
    await page.locator("#search-dialog input[type='search']").fill("config");
    await expect(page.locator(".search-result").first()).toBeVisible({
      timeout: 5000,
    });
    await page.locator(".search-result").first().click();
    await expect(page).toHaveURL(/\/guide\/config\//);
  });

  test("shows no results for gibberish query", async ({ page }) => {
    await page.locator('button[aria-label="Search"]').click();
    await page
      .locator("#search-dialog input[type='search']")
      .fill("xyzzyplugh");
    await expect(page.locator("[data-search-status]")).toContainText(
      "No results",
      { timeout: 5000 }
    );
    expect(await page.locator(".search-result").count()).toBe(0);
  });

  test("clears results when dialog is reopened", async ({ page }) => {
    await page.locator('button[aria-label="Search"]').click();
    await page.locator("#search-dialog input[type='search']").fill("config");
    await expect(page.locator(".search-result").first()).toBeVisible({
      timeout: 5000,
    });
    await page.keyboard.press("Escape");
    await expect(page.locator("#search-dialog")).not.toBeVisible();
    await page.locator('button[aria-label="Search"]').click();
    await expect(
      page.locator("#search-dialog input[type='search']")
    ).toHaveValue("");
    expect(await page.locator(".search-result").count()).toBe(0);
  });

  test("scores title matches higher than body matches", async ({ page }) => {
    await page.locator('button[aria-label="Search"]').click();
    await page
      .locator("#search-dialog input[type='search']")
      .fill("configuration");
    await expect(page.locator(".search-result").first()).toBeVisible({
      timeout: 5000,
    });
    await expect(
      page.locator(".search-result").first().locator("strong")
    ).toContainText("Configuration");
  });
});
