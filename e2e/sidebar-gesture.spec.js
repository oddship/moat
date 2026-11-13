const { test, expect } = require("./fixtures");

async function swipe(page, target, dx, dy = 0) {
  const box = await target.boundingBox();
  if (!box) throw new Error("Swipe target is not visible");
  const start = { x: box.x + 24, y: box.y + box.height / 2 };
  const client = await page.context().newCDPSession(page);
  try {
    await client.send("Emulation.setTouchEmulationEnabled", { enabled: true });
    await client.send("Input.dispatchTouchEvent", { type: "touchStart", touchPoints: [{ ...start, id: 1 }] });
    for (let step = 1; step <= 5; step++) {
      await client.send("Input.dispatchTouchEvent", {
        type: "touchMove",
        touchPoints: [{ x: start.x + dx * step / 5, y: start.y + dy * step / 5, id: 1 }],
      });
    }
    await client.send("Input.dispatchTouchEvent", { type: "touchEnd", touchPoints: [] });
  } finally {
    await client.detach();
  }
}

test.describe("Mobile sidebar swipe", () => {
  test.beforeEach(async ({ page }) => {
    await page.setViewportSize({ width: 390, height: 844 });
    await page.goto("/guide/quickstart/");
  });

  test("a right swipe opens the sidebar and keeps its toggle synchronized", async ({ page }) => {
    const toggle = page.locator("[data-sidebar-toggle]");
    await swipe(page, page.locator("article p").first(), 140);
    await expect(page.locator("body")).toHaveAttribute("data-sidebar-open", "");
    await expect(toggle).toHaveAttribute("aria-expanded", "true");
    await expect(page.locator("aside[data-sidebar]")).toBeVisible();
    await page.keyboard.press("Escape");
    await expect(toggle).toHaveAttribute("aria-expanded", "false");
  });

  test("vertical scrolling and short horizontal movements keep the sidebar closed", async ({ page }) => {
    const paragraph = page.locator("article p").first();
    await swipe(page, paragraph, 5, -100);
    await expect(page.locator("[data-sidebar-toggle]")).toHaveAttribute("aria-expanded", "false");
    await paragraph.scrollIntoViewIfNeeded();
    await swipe(page, paragraph, 25);
    await expect(page.locator("[data-sidebar-toggle]")).toHaveAttribute("aria-expanded", "false");
  });

  test("swiping a code block does not open the sidebar", async ({ page }) => {
    await swipe(page, page.locator("pre").first(), 140);
    await expect(page.locator("[data-sidebar-toggle]")).toHaveAttribute("aria-expanded", "false");
  });

  test("a horizontally scrollable region keeps its normal touch behavior", async ({ page }) => {
    await page.evaluate(() => {
      const scroller = document.createElement("div");
      scroller.id = "gesture-test-scroller";
      scroller.style.cssText = "position:fixed;left:20px;top:420px;width:240px;height:60px;overflow-x:auto";
      scroller.innerHTML = '<div style="width:800px;height:60px">Scrollable content</div>';
      document.body.append(scroller);
      scroller.scrollLeft = 200;
    });
    const scroller = page.locator("#gesture-test-scroller");
    await swipe(page, scroller, 140);
    await expect(page.locator("[data-sidebar-toggle]")).toHaveAttribute("aria-expanded", "false");
    await expect.poll(() => scroller.evaluate(element => element.scrollLeft)).toBeLessThan(200);
  });

  test("touch swipes do not toggle the desktop sidebar", async ({ page }) => {
    await page.setViewportSize({ width: 1280, height: 900 });
    await swipe(page, page.locator("article p").first(), 140);
    await expect(page.locator("body")).not.toHaveAttribute("data-sidebar-open", "");
  });
});
