const { test, expect } = require("./fixtures");

test.describe("Agent-friendly pages", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto("/guide/quickstart/");
  });

  test("copies the exported Markdown without the site layout", async ({ page }) => {
    // Capture the clipboard payload without relying on the shared browser's
    // permission settings. The click and Markdown fetch still run normally.
    await page.evaluate(() => {
      Object.defineProperty(navigator, "clipboard", {
        configurable: true,
        value: {
          writeText: async (text) => { window.copiedPageMarkdown = text; },
        },
      });
    });
    const markdownURL = await page.getByRole("link", { name: "View Markdown", exact: true }).getAttribute("href");
    const response = await page.request.get(new URL(markdownURL, page.url()).href);
    expect(response.ok()).toBe(true);
    const markdown = await response.text();
    await page.getByRole("button", { name: "Copy page", exact: true }).click();
    await expect(page.locator("[data-copy-status]")).toHaveText("Copied!");
    expect(await page.evaluate(() => window.copiedPageMarkdown)).toBe(markdown);
    expect(markdown).toContain("# Quick start");
    expect(markdown).not.toContain("data-topnav");
  });

  test("keeps a manual fallback when the Markdown request fails", async ({ page }) => {
    await page.route("**/guide/quickstart/index.md", route => route.fulfill({ status: 503, body: "Unavailable" }));
    await page.getByRole("button", { name: "Copy page", exact: true }).click();
    await expect(page.locator("[data-copy-status]")).toContainText("Use View Markdown");
    await expect(page.getByRole("button", { name: "Copy page", exact: true })).toBeEnabled();
    await expect(page.getByRole("link", { name: "View Markdown", exact: true })).toBeVisible();
  });

  test("opens the page export through View Markdown", async ({ page }) => {
    await page.getByRole("link", { name: "View Markdown", exact: true }).click();
    await expect(page).toHaveURL(/\/guide\/quickstart\/index\.md$/);
    await expect(page.locator("body")).toContainText("# Quick start");
  });
});
