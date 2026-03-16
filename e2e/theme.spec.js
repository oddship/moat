const { test, expect } = require("./fixtures");

test.describe("Theme toggle", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto("/");
    await page.evaluate(() => localStorage.removeItem("theme"));
  });

  test("toggles theme on button click", async ({ page }) => {
    const html = page.locator("html").first();
    const initialTheme = await html.getAttribute("data-theme");
    await page.locator("aside[data-sidebar] footer button").click();
    const newTheme = await html.getAttribute("data-theme");
    expect(newTheme).not.toEqual(initialTheme);
  });

  test("persists theme in localStorage", async ({ page }) => {
    await page.locator("aside[data-sidebar] footer button").click();
    const stored = await page.evaluate(() => localStorage.getItem("theme"));
    expect(stored).toBeTruthy();
  });

  test("restores theme on page reload", async ({ page }) => {
    await page.locator("aside[data-sidebar] footer button").click();
    const theme = await page.locator("html").first().getAttribute("data-theme");
    await page.reload();
    await expect(page.locator("html").first()).toHaveAttribute(
      "data-theme",
      theme
    );
  });
});
