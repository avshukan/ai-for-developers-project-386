/**
 * Date/time helpers.
 *
 * The API uses UTC ISO 8601 strings (see docs/domain.md "Time Rules").
 * The UI displays them in the browser's local time zone.
 */

const dateTimeFormatter = new Intl.DateTimeFormat(undefined, {
  weekday: 'short',
  day: 'numeric',
  month: 'short',
  hour: '2-digit',
  minute: '2-digit',
});

const timeFormatter = new Intl.DateTimeFormat(undefined, {
  hour: '2-digit',
  minute: '2-digit',
});

const dayFormatter = new Intl.DateTimeFormat(undefined, {
  weekday: 'long',
  day: 'numeric',
  month: 'long',
});

/** "Mon, 24 Jun, 10:00" — full local date and time. */
export function formatDateTime(iso: string): string {
  const date = new Date(iso);
  return Number.isNaN(date.getTime()) ? iso : dateTimeFormatter.format(date);
}

/** "10:00" — local time only (used for slot buttons grouped under a day). */
export function formatTime(iso: string): string {
  const date = new Date(iso);
  return Number.isNaN(date.getTime()) ? iso : timeFormatter.format(date);
}

/** "Monday, 24 June" — local day heading (used to group slots). */
export function formatDay(iso: string): string {
  const date = new Date(iso);
  return Number.isNaN(date.getTime()) ? iso : dayFormatter.format(date);
}

/** Stable per-local-day key for grouping slots, e.g. "2026-5-24". */
export function dayKey(iso: string): string {
  const date = new Date(iso);
  if (Number.isNaN(date.getTime())) return iso;
  return `${date.getFullYear()}-${date.getMonth()}-${date.getDate()}`;
}
