import { test, expect, type Page } from '@playwright/test';

/**
 * Integration coverage for the main booking scenario (docs/product.md →
 * "Main Guest Scenario" and "Main Host Bookings Scenario"). These tests drive
 * the real SPA against the real backend; nothing is mocked. The backend serves
 * its demo seed (two event types + availability), so available slots exist.
 *
 * The backend store is in memory and shared across the run, so each test books
 * its own *first available* slot rather than a fixed one, and assertions are
 * written to tolerate slots removed by earlier tests.
 */

const EVENT_TYPE = /Intro call/;

/** A unique guest per test so host-side assertions can target one booking. */
function uniqueGuest() {
  const tag = `${Date.now()}-${Math.floor(Math.random() * 1e6)}`;
  return { name: `Guest ${tag}`, email: `guest-${tag}@example.com` };
}

/** Open the booking page and choose the event type under test. */
async function chooseEventType(page: Page) {
  await page.goto('/');
  await expect(page.getByRole('heading', { name: 'Book a meeting' })).toBeVisible();
  // Guard: the SPA must be wired to the real backend, not the Prism mock default,
  // otherwise these tests would pass against canned example data.
  await expect(page.getByText('API: http://localhost:8080')).toBeVisible();
  await page.getByRole('button', { name: EVENT_TYPE }).click();
  // Wait for the slots section to render.
  await expect(page.getByRole('heading', { name: '2. Pick a slot' })).toBeVisible();
}

/** Capture and click the first available slot. Returns its day + time labels. */
async function pickFirstSlot(page: Page): Promise<{ day: string; time: string }> {
  const firstDay = page.locator('.slots__day').first();
  await expect(firstDay).toBeVisible();
  const day = (await firstDay.locator('.slots__day-title').innerText()).trim();
  const firstSlot = firstDay.locator('.slots__list button').first();
  const time = (await firstSlot.innerText()).trim();
  await firstSlot.click();
  return { day, time };
}

test('guest books an available slot end to end', async ({ page }) => {
  const guest = uniqueGuest();

  await chooseEventType(page);
  await pickFirstSlot(page);

  // Step 3 — guest details.
  await expect(page.getByRole('heading', { name: '3. Your details' })).toBeVisible();
  await page.getByLabel('Name', { exact: true }).fill(guest.name);
  await page.getByLabel('Email', { exact: true }).fill(guest.email);
  await page.getByRole('button', { name: /confirm booking/i }).click();

  // Confirmation reflects what the backend created.
  await expect(page.getByRole('heading', { name: 'Booking confirmed' })).toBeVisible();
  await expect(page.getByText(`${guest.name} (${guest.email})`)).toBeVisible();
});

test('a booked slot becomes unavailable (no double booking)', async ({ page }) => {
  const guest = uniqueGuest();

  await chooseEventType(page);
  const { day, time } = await pickFirstSlot(page);

  await page.getByLabel('Name', { exact: true }).fill(guest.name);
  await page.getByLabel('Email', { exact: true }).fill(guest.email);
  await page.getByRole('button', { name: /confirm booking/i }).click();
  await expect(page.getByRole('heading', { name: 'Booking confirmed' })).toBeVisible();

  // Reload the same event type's slots and confirm the booked slot is gone.
  await page.getByRole('button', { name: /book another/i }).click();
  await expect(page.getByRole('heading', { name: '2. Pick a slot' })).toBeVisible();

  const sameDay = page.locator('.slots__day').filter({ hasText: day });
  await expect(
    sameDay.locator('.slots__list button', { hasText: new RegExp(`^${time}$`) }),
  ).toHaveCount(0);
});

test('host sees a new booking on the bookings page', async ({ page }) => {
  const guest = uniqueGuest();

  // Guest books a slot.
  await chooseEventType(page);
  await pickFirstSlot(page);
  await page.getByLabel('Name', { exact: true }).fill(guest.name);
  await page.getByLabel('Email', { exact: true }).fill(guest.email);
  await page.getByRole('button', { name: /confirm booking/i }).click();
  await expect(page.getByRole('heading', { name: 'Booking confirmed' })).toBeVisible();

  // Host opens the bookings page and finds it.
  await page.goto('/host/bookings');
  await expect(page.getByRole('heading', { name: 'Upcoming bookings' })).toBeVisible();
  await expect(page.getByText(guest.email)).toBeVisible();
});
