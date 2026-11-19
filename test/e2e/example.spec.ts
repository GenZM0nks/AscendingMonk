import { test, expect } from "@playwright/test";

test("can log in user", async ({ page }) => {
  await page.goto("/login");
  await expect(page).toHaveTitle("Login · ¿Who Knows?");

  let navLinks = page.getByRole("navigation").getByRole("link");
  await expect(navLinks).toHaveText(["¿Who Knows?", "Log in", "Register", "Weather"]);

  await page.locator('input[name="username"]').fill(process.env.TEST_USERNAME!);
  await page.locator('input[name="password"]').fill(process.env.TEST_PASSWORD!);
  await page.getByRole("button", { name: "Log in" }).click();
  await expect(page).toHaveURL("/");

  navLinks = page.getByRole("navigation").getByRole("link");
  await expect(navLinks).toHaveText(["¿Who Knows?", "Log out", "Weather"]);
});

test("can search", async ({ page }) => {
  await page.goto("/");
  await expect(page).toHaveTitle("Search · ¿Who Knows?");
  await page.getByRole("textbox").fill("Fortran");
  await page.getByRole("button", { name: "Search" }).click();
  const firstResult = page.locator('#results a.search-result-title').first();
  await expect(firstResult).toContainText('Fortran');
});
