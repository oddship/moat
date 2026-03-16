const { defineConfig } = require("@playwright/test");

module.exports = defineConfig({
  testDir: "./e2e",
  timeout: 15000,
  workers: 1,
  use: {
    baseURL: "http://localhost:8080",
  },
});
