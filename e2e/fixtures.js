const base = require("@playwright/test");
const { chromium } = require("playwright-core");

const BASE_URL = "http://localhost:8080";

// Custom fixture that connects to existing Chromium via CDP
// instead of launching a new browser.
const test = base.test.extend({
  context: async ({}, use) => {
    const browser = await chromium.connectOverCDP("http://localhost:9222");
    const context =
      browser.contexts()[0] || (await browser.newContext({ baseURL: BASE_URL }));
    await use(context);
  },

  page: async ({ context }, use) => {
    const page = await context.newPage();
    // Manually set baseURL for goto("/") to work
    const originalGoto = page.goto.bind(page);
    page.goto = async (url, options) => {
      if (url.startsWith("/")) {
        url = BASE_URL + url;
      }
      return originalGoto(url, options);
    };
    await use(page);
    await page.close();
  },
});

module.exports = { test, expect: base.expect };
