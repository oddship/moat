const { test, expect } = require("./fixtures");

test.describe("Layout", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto("/");
  });

  test("renders the top nav with branding", async ({ page }) => {
    const topnav = page.locator("nav[data-topnav]");
    await expect(topnav).toBeVisible();
    const brandLink = topnav.locator("a").first();
    await expect(brandLink).toHaveAttribute("href", "/");
  });

  test("renders topnav GitHub link with icon", async ({ page }) => {
    const ghLink = page.locator("nav[data-topnav] .col-8 a");
    await expect(ghLink).toContainText("GitHub");
    await expect(ghLink).toHaveAttribute(
      "href",
      "https://github.com/oddship/moat"
    );
    await expect(ghLink.locator("svg")).toBeAttached();
  });

  test("renders the sidebar with navigation", async ({ page }) => {
    const sidebar = page.locator("aside[data-sidebar]");
    await expect(sidebar).toBeAttached();
    const items = sidebar.locator("nav ul li");
    expect(await items.count()).toBeGreaterThanOrEqual(3);
  });

  test("renders sidebar GitHub link with icon", async ({ page }) => {
    const ghLink = page.locator("aside[data-sidebar] nav a").first();
    await expect(ghLink).toContainText("GitHub");
    await expect(ghLink.locator("svg")).toBeAttached();
  });

  test("renders collapsible nav sections", async ({ page }) => {
    const sections = page.locator("aside[data-sidebar] details");
    expect(await sections.count()).toBeGreaterThanOrEqual(2);
    await expect(sections.first().locator("summary")).toContainText("Guide");
  });

  test("renders the theme toggle button", async ({ page }) => {
    await expect(
      page.locator("aside[data-sidebar] footer button")
    ).toContainText("Switch theme");
  });

  test("has no inline style attributes in topnav or main", async ({
    page,
  }) => {
    const topnavStyles = await page
      .locator("nav[data-topnav] [style]")
      .count();
    expect(topnavStyles).toBe(0);
    const mainStyles = await page.locator("main [style]").count();
    expect(mainStyles).toBe(0);
  });

  test("navigates to an inner page", async ({ page }) => {
    await page.locator("aside[data-sidebar] a", { hasText: "Configuration" }).click();
    await expect(page).toHaveURL(/\/guide\/config\//);
    await expect(page.locator("article h1")).toContainText("Configuration");
  });

  test("marks the current page in sidebar nav", async ({ page }) => {
    await page.goto("/guide/config/");
    const active = page.locator('aside[data-sidebar] a[aria-current="page"]');
    await expect(active).toContainText("Configuration");
  });
});
